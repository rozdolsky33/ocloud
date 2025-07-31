package display

import (
	"fmt"
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

	fmt.Println()
}

// displayBanner displays the OCloud ASCII art banner
func displayBanner() {
	fmt.Println(" ██████╗  ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ")
	fmt.Println("██╔═══██╗██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗")
	fmt.Println("██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║")
	fmt.Println("██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║")
	fmt.Println("╚██████╔╝╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝")
	fmt.Println(" ╚═════╝  ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝")
	fmt.Println()
}
