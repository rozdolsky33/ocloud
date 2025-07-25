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
func AuthenticateWithOCI(filter, realm string) error {

	s := NewService()

	logger.LogWithLevel(s.logger, 1, "Authenticating with OCI", "filter", filter, "realm", realm)

	result, err = performInteractiveAuthentication(s, filter, realm)
	if err != nil {
		return fmt.Errorf("performing interactive authentication: %w", err)
	}

	logger.LogWithLevel(s.logger, 3, "Interactive authentication completed", "tenancyID", result.TenancyID, "tenancyName", result.TenancyName)

	// Display environment variables
	logger.LogWithLevel(s.logger, 3, "Displaying environment variables")
	err = PrintExportVariable(result.Profile, result.TenancyName, result.CompartmentName)

	if err != nil {
		return fmt.Errorf("printing export variables: %w", err)
	}

	logger.LogWithLevel(s.logger, 1, "Authentication process completed successfully")
	return nil
}

// performInteractiveAuthentication handles the interactive authentication process.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
func performInteractiveAuthentication(s *Service, filter, realm string) (*AuthenticationResult, error) {
	// Prompt for profile selection
	profile, err := s.PromptForProfile()
	if err != nil {
		return nil, fmt.Errorf("selecting profile: %w", err)
	}
	logger.LogWithLevel(s.logger, 1, "Profile selected", "profile", profile)

	// Display regions in a table
	logger.LogWithLevel(s.logger, 3, "Getting OCI regions")
	regions := s.GetOCIRegions()
	logger.LogWithLevel(s.logger, 3, "Displaying regions table", "regionCount", len(regions), "filter", filter)

	if err := DisplayRegionsTable(regions, filter); err != nil {
		return nil, fmt.Errorf("displaying regions: %w", err)
	}

	// Prompt for region selection
	region, err := s.PromptForRegion()
	if err != nil {
		return nil, fmt.Errorf("selecting region: %w", err)
	}

	fmt.Printf("Using region: %s\n", region)
	logger.LogWithLevel(s.logger, 3, "Region selected", "region", region)

	// Authenticate with OCI
	logger.LogWithLevel(s.logger, 3, "Authenticating with OCI", "profile", profile, "region", region)
	result, err := s.Authenticate(profile, region)

	if err != nil {
		return nil, fmt.Errorf("authenticating with OCI: %w", err)
	}

	logger.LogWithLevel(s.logger, 3, "Authentication successful", "profile", profile, "region", region)

	// View configuration
	err = info.ViewConfiguration(false, realm)
	if err != nil {
		return nil, fmt.Errorf("viewing configuration: %w", err)
	}

	// Prompt for custom environment variables
	if s.promptYesNo("Do you want to set OCI_TENANCY_NAME and OCI_COMPARTMENT?") {
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

		if tenancy != "" {
			result.TenancyName = tenancy
			logger.LogWithLevel(s.logger, 3, "Updated tenancy name", "tenancyName", tenancy)
		}

		if compartment != "" {
			result.CompartmentName = compartment
			logger.LogWithLevel(s.logger, 3, "Updated compartment", "compartment", compartment)
		}

	} else {
		logger.LogWithLevel(s.logger, 3, "Skipping variable setup")
		fmt.Println("\n Skipping variable setup.")
	}

	logger.LogWithLevel(s.logger, 1, "Interactive authentication completed successfully", "profile", profile, "region", region)
	return result, nil
}
