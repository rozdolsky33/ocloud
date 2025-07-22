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
// If envOnly is true, it skips the interactive prompts and only outputs the environment variables.
// If the filter is not empty, it filters the regions by prefix.
func AuthenticateWithOCI(appCtx *app.ApplicationContext, envOnly bool, filter string) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Authenticating with OCI", "envOnly", envOnly)

	// Create a new service
	service := NewService(appCtx)

	if envOnly {
		// In env-only mode, skip interactive prompts and authentication
		// Just get the current environment variables
		_, err := service.GetCurrentEnvironment()
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
		_, err = service.Authenticate(profile, region)
		if err != nil {
			return fmt.Errorf("authenticating with OCI: %w", err)
		}
		err = info.ViewConfiguration(appCtx, false, "")
		if err != nil {
			return fmt.Errorf("viewing configuration: %w", err)
		}

		err = PrintExportVariable()
		fmt.Println("\nExport the following environment variables to use them in the future")

		if promptYesNo("Do you want to set OCI_TENANCY_NAME and OCI_COMPARTMENT?") {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter OCI_TENANCY_NAME: ")
			tenancy, _ := reader.ReadString('\n')

			fmt.Print("Enter OCI_COMPARTMENT: ")
			compartment, _ := reader.ReadString('\n')

			tenancy = strings.TrimSpace(tenancy)
			compartment = strings.TrimSpace(compartment)

			fmt.Println()
			// Output export statements
			fmt.Printf("export OCI_TENANCY_NAME=%q\n", tenancy)
			fmt.Printf("export OCI_COMPARTMENT=%q\n", compartment)
		} else {
			fmt.Println("Skipping variable setup.")
		}

	}

	return nil
}
