package auth

import (
	"bufio"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/info"
	"os"
	"strings"
)

// AuthenticateWithOCI handles the authentication process with OCI.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
// If the filter is not empty, it filters the regions by prefix.
func AuthenticateWithOCI(appCtx *app.ApplicationContext, filter string) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Authenticating with OCI", "filter", filter)

	// Create a new service
	logger.LogWithLevel(appCtx.Logger, 3, "Creating authentication service")
	service := NewService(appCtx)

	var result *AuthenticationResult
	var err error

	// Perform interactive authentication
	logger.LogWithLevel(appCtx.Logger, 3, "Starting interactive authentication")
	result, err = performInteractiveAuthentication(appCtx, service, filter)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to perform interactive authentication", "error", err)
		return fmt.Errorf("performing interactive authentication: %w", err)
	}
	logger.LogWithLevel(appCtx.Logger, 3, "Interactive authentication completed", "tenancyID", result.TenancyID, "tenancyName", result.TenancyName)

	// Display environment variables
	logger.LogWithLevel(appCtx.Logger, 3, "Displaying environment variables")
	err = PrintExportVariable(result.TenancyName, result.TenancyID)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to print export variables", "error", err)
		return fmt.Errorf("printing export variables: %w", err)
	}

	logger.LogWithLevel(appCtx.Logger, 3, "Displaying instructions for persisting environment variables")
	fmt.Println("\nTo persist your selection, export the following environment variables in your shell:")

	logger.LogWithLevel(appCtx.Logger, 1, "Authentication process completed successfully")
	return nil
}

// performInteractiveAuthentication handles the interactive authentication process.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
func performInteractiveAuthentication(appCtx *app.ApplicationContext, service *Service, filter string) (*AuthenticationResult, error) {
	logger.LogWithLevel(appCtx.Logger, 1, "Starting interactive authentication process")

	// Prompt for profile selection
	logger.LogWithLevel(appCtx.Logger, 3, "Prompting for profile selection")
	profile, err := service.PromptForProfile()
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to select profile", "error", err)
		return nil, fmt.Errorf("selecting profile: %w", err)
	}
	logger.LogWithLevel(appCtx.Logger, 3, "Profile selected", "profile", profile)

	// Display regions in a table
	logger.LogWithLevel(appCtx.Logger, 3, "Getting OCI regions")
	regions := service.GetOCIRegions()
	logger.LogWithLevel(appCtx.Logger, 3, "Displaying regions table", "regionCount", len(regions), "filter", filter)
	if err := DisplayRegionsTable(regions, appCtx, filter); err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to display regions", "error", err)
		return nil, fmt.Errorf("displaying regions: %w", err)
	}

	// Prompt for region selection
	logger.LogWithLevel(appCtx.Logger, 3, "Prompting for region selection")
	region, err := service.PromptForRegion()
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to select region", "error", err)
		return nil, fmt.Errorf("selecting region: %w", err)
	}
	fmt.Printf("Using region: %s\n", region)
	logger.LogWithLevel(appCtx.Logger, 3, "Region selected", "region", region)

	// Authenticate with OCI
	logger.LogWithLevel(appCtx.Logger, 3, "Authenticating with OCI", "profile", profile, "region", region)
	result, err := service.Authenticate(profile, region)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to authenticate with OCI", "error", err)
		return nil, fmt.Errorf("authenticating with OCI: %w", err)
	}
	logger.LogWithLevel(appCtx.Logger, 3, "Authentication successful", "profile", profile, "region", region)

	// View configuration
	logger.LogWithLevel(appCtx.Logger, 3, "Viewing configuration")
	err = info.ViewConfiguration(appCtx, false, "")
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, 1, "Failed to view configuration", "error", err)
		return nil, fmt.Errorf("viewing configuration: %w", err)
	}

	// Prompt for custom environment variables
	logger.LogWithLevel(appCtx.Logger, 3, "Prompting for custom environment variables")
	if promptYesNo("Do you want to set OCI_TENANCY_NAME and OCI_COMPARTMENT?") {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter OCI_TENANCY_NAME: ")
		tenancy, err := reader.ReadString('\n')
		if err != nil {
			logger.LogWithLevel(appCtx.Logger, 3, "Error reading tenancy name input", "error", err)
		}

		fmt.Print("Enter OCI_COMPARTMENT: ")
		compartment, err := reader.ReadString('\n')
		if err != nil {
			logger.LogWithLevel(appCtx.Logger, 3, "Error reading compartment input", "error", err)
		}

		tenancy = strings.TrimSpace(tenancy)
		compartment = strings.TrimSpace(compartment)
		logger.LogWithLevel(appCtx.Logger, 3, "Custom environment variables entered", "tenancyName", tenancy, "compartment", compartment)

		// Update the result with custom values
		if tenancy != "" {
			result.TenancyName = tenancy
			logger.LogWithLevel(appCtx.Logger, 3, "Updated tenancy name", "tenancyName", tenancy)
		}

		// Store compartment in TenancyID for now, as we don't have a separate field for it
		if compartment != "" {
			result.TenancyID = compartment
			logger.LogWithLevel(appCtx.Logger, 3, "Updated compartment", "compartment", compartment)
		}
	} else {
		logger.LogWithLevel(appCtx.Logger, 3, "Skipping variable setup")
		fmt.Println("Skipping variable setup.")
	}

	logger.LogWithLevel(appCtx.Logger, 1, "Interactive authentication completed successfully", "profile", profile, "region", region)
	return result, nil
}
