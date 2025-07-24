package auth

import (
	"bufio"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/info"
	"os"
	"strings"
)

// AuthenticateWithOCI handles the authentication process with OCI.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
// If the filter is not empty, it filters the regions by prefix.
func AuthenticateWithOCI(filter string) error {

	s := NewService()

	logger.LogWithLevel(s.logger, 1, "Authenticating with OCI", "filter", filter)

	// Create a new service
	logger.LogWithLevel(s.logger, 3, "Creating authentication service")

	var result *AuthenticationResult
	var err error

	// Perform interactive authentication
	logger.LogWithLevel(s.logger, 3, "Starting interactive authentication")
	result, err = performInteractiveAuthentication(s, filter)
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to perform interactive authentication", "error", err)
		return fmt.Errorf("performing interactive authentication: %w", err)
	}
	logger.LogWithLevel(s.logger, 3, "Interactive authentication completed", "tenancyID", result.TenancyID, "tenancyName", result.TenancyName)

	// Display environment variables
	logger.LogWithLevel(s.logger, 3, "Displaying environment variables")
	err = PrintExportVariable(result.TenancyName, result.TenancyID)
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to print export variables", "error", err)
		return fmt.Errorf("printing export variables: %w", err)
	}

	logger.LogWithLevel(s.logger, 3, "Displaying instructions for persisting environment variables")
	fmt.Println("\nTo persist your selection, export the following environment variables in your shell:")

	logger.LogWithLevel(s.logger, 1, "Authentication process completed successfully")
	return nil
}

// performInteractiveAuthentication handles the interactive authentication process.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
func performInteractiveAuthentication(s *Service, filter string) (*AuthenticationResult, error) {
	logger.LogWithLevel(s.logger, 1, "Starting interactive authentication process")

	// Prompt for profile selection
	logger.LogWithLevel(s.logger, 3, "Prompting for profile selection")
	profile, err := s.PromptForProfile()
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to select profile", "error", err)
		return nil, fmt.Errorf("selecting profile: %w", err)
	}
	logger.LogWithLevel(s.logger, 3, "Profile selected", "profile", profile)

	// Display regions in a table
	logger.LogWithLevel(s.logger, 3, "Getting OCI regions")
	regions := s.GetOCIRegions()
	logger.LogWithLevel(s.logger, 3, "Displaying regions table", "regionCount", len(regions), "filter", filter)
	if err := DisplayRegionsTable(regions, filter); err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to display regions", "error", err)
		return nil, fmt.Errorf("displaying regions: %w", err)
	}

	// Prompt for region selection
	logger.LogWithLevel(s.logger, 3, "Prompting for region selection")
	region, err := s.PromptForRegion()
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to select region", "error", err)
		return nil, fmt.Errorf("selecting region: %w", err)
	}
	fmt.Printf("Using region: %s\n", region)
	logger.LogWithLevel(s.logger, 3, "Region selected", "region", region)

	// Authenticate with OCI
	logger.LogWithLevel(s.logger, 3, "Authenticating with OCI", "profile", profile, "region", region)
	result, err := s.Authenticate(profile, region)
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to authenticate with OCI", "error", err)
		return nil, fmt.Errorf("authenticating with OCI: %w", err)
	}
	logger.LogWithLevel(s.logger, 3, "Authentication successful", "profile", profile, "region", region)

	// View configuration
	logger.LogWithLevel(s.logger, 3, "Viewing configuration")
	err = info.ViewConfiguration(false, "")
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "Failed to view configuration", "error", err)
		return nil, fmt.Errorf("viewing configuration: %w", err)
	}

	// Prompt for custom environment variables
	logger.LogWithLevel(s.logger, 3, "Prompting for custom environment variables")
	if promptYesNo("Do you want to set OCI_TENANCY_NAME and OCI_COMPARTMENT?") {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter OCI_TENANCY_NAME: ")
		tenancy, err := reader.ReadString('\n')
		if err != nil {
			logger.LogWithLevel(s.logger, 3, "Error reading tenancy name input", "error", err)
		}

		fmt.Print("Enter OCI_COMPARTMENT: ")
		compartment, err := reader.ReadString('\n')
		if err != nil {
			logger.LogWithLevel(s.logger, 3, "Error reading compartment input", "error", err)
		}

		tenancy = strings.TrimSpace(tenancy)
		compartment = strings.TrimSpace(compartment)
		logger.LogWithLevel(s.logger, 3, "Custom environment variables entered", "tenancyName", tenancy, "compartment", compartment)

		// Update the result with custom values
		if tenancy != "" {
			result.TenancyName = tenancy
			logger.LogWithLevel(s.logger, 3, "Updated tenancy name", "tenancyName", tenancy)
		}

		// Store compartment in TenancyID for now, as we don't have a separate field for it
		if compartment != "" {
			result.TenancyID = compartment
			logger.LogWithLevel(s.logger, 3, "Updated compartment", "compartment", compartment)
		}
	} else {
		logger.LogWithLevel(s.logger, 3, "Skipping variable setup")
		fmt.Println("Skipping variable setup.")
	}

	logger.LogWithLevel(s.logger, 1, "Interactive authentication completed successfully", "profile", profile, "region", region)
	return result, nil
}
