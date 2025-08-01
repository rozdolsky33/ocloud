package display

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/buildinfo"
	"github.com/rozdolsky33/ocloud/internal/config"
	"os"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

// PrintOCIConfiguration displays the current configuration details
// and checks if required environment variables are set
func PrintOCIConfiguration() {
	displayBanner()

	fmt.Println("\033[1mConfiguration Details:\033[0m")

	profile := os.Getenv("OCI_CLI_PROFILE")
	if profile == "" {
		fmt.Println("  \033[33mOCI_CLI_PROFILE\033[0m: \033[31mNot set - Please set profile\033[0m")
	} else {
		fmt.Printf("  \033[33mOCI_CLI_PROFILE\033[0m: %s\n", profile)
	}

	tenancyName := os.Getenv(flags.EnvOCITenancyName)
	if tenancyName == "" {
		fmt.Println("  \033[33mOCI_TENANCY_NAME\033[0m: \033[31mNot set - Please set tenancy\033[0m")
	} else {
		fmt.Printf("  \033[33mOCI_TENANCY_NAME\033[0m: %s\n", tenancyName)
	}

	compartment := os.Getenv(flags.EnvOCICompartment)
	if compartment == "" {
		fmt.Println("  \033[33mOCI_COMPARTMENT\033[0m: \033[31mNot set - Please set compartmen name\033[0m")
	} else {
		fmt.Printf("  \033[33mOCI_COMPARTMENT\033[0m: %s\n", compartment)
	}
	path := config.TenancyMapPath()
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		// This block executes if the file does not exist.
		fmt.Println("  \033[33mOCI_TENANCY_MAP_PATH\033[0m: \033[31mNot set (file not found)\033[0m")
	} else if err != nil {
		// This block handles other potential errors, e.g., permission denied.
		fmt.Printf("  \033[33mOCI_TENANCY_MAP_PATH\033[0m: \033[31mError checking file: %v\033[0m\n", err)
	} else {
		// If err is nil, the stat was successful and the file exists.
		fmt.Printf("  \033[33mOCI_TENANCY_MAP_PATH\033[0m: %s\n", path)
	}

	fmt.Println()
}

// displayBanner displays the OCloud ASCII art banner
func displayBanner() {
	fmt.Println()
	fmt.Println(" ██████╗  ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ")
	fmt.Println("██╔═══██╗██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗")
	fmt.Println("██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║")
	fmt.Println("██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║")
	fmt.Println("╚██████╔╝╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝")
	fmt.Println(" ╚═════╝  ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝")
	fmt.Println()
	fmt.Printf("  \033[33mVersion\033[0m: \033[32m%s\033[0m\n", buildinfo.Version)
	fmt.Println()
}
