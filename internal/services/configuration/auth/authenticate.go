package auth

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// AuthenticateWithOCI handles the authentication process with OCI.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
// If envOnly is true, it skips the interactive prompts and only outputs the environment variables.
// If the filter is not empty, it filters the regions by prefix.
func AuthenticateWithOCI(appCtx *app.ApplicationContext, envOnly bool, filter string) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Authenticating with OCI", "envOnly", envOnly)

	// Create a new service
	service := NewService(appCtx)

	var result *AuthenticationResult
	var err error

	if envOnly {
		// In env-only mode, skip interactive prompts and authentication
		// Just get the current environment variables
		result, err = service.GetCurrentEnvironment()
		if err != nil {
			return fmt.Errorf("getting current environment: %w", err)
		}
	} else {
		// Prompt for profile selection
		profile, err := service.PromptForProfile()
		if err != nil {
			return fmt.Errorf("selecting profile: %w", err)
		}

		// Display regions in a table
		regions := service.GetOCIRegions()
		if err := DisplayRegionsTable(regions, appCtx, filter); err != nil {
			return fmt.Errorf("displaying regions: %w", err)
		}

		// Prompt for region selection
		region, err := service.PromptForRegion()
		if err != nil {
			return fmt.Errorf("selecting region: %w", err)
		}
		fmt.Printf("Using region: %s\n", region)

		// Authenticate with OCI
		result, err = service.Authenticate(profile, region)
		if err != nil {
			return fmt.Errorf("authenticating with OCI: %w", err)
		}
	}

	// Print export for compartment and tenancy
	fmt.Printf("export OCI_COMPARTMENT=%s\n", result.TenancyID)
	fmt.Printf("export OCI_CLI_TENANCY=%s\n", result.TenancyID)

	// Print export for tenancy name if available
	if result.TenancyName != "" {
		fmt.Printf("export OCI_TENANCY_NAME=%s\n", result.TenancyName)
	}

	if !envOnly {
		fmt.Println("✅ Authentication complete. Run `eval $(ocloud config auth -e)` to set your env.")
	}

	return nil
}
