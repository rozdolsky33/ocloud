package display

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/buildinfo"
	"github.com/rozdolsky33/ocloud/internal/config"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

var validRe = regexp.MustCompile(`(?i)^Session is valid until\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s*$`)
var expiredRe = regexp.MustCompile(`(?i)^Session has expired\s*$`)

// CheckOCISessionValidity checks the validity of the OCI session
func CheckOCISessionValidity() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "oci", "session", "validate", "--local")
	out, err := cmd.CombinedOutput()
	raw := strings.TrimSpace(string(out))

	if matches := validRe.FindStringSubmatch(raw); len(matches) > 1 {
		return fmt.Sprintf("\033[32mValid until %s\033[0m", matches[1])
	} else if expiredRe.MatchString(raw) {
		return "\033[31mSession Expired\033[0m"
	} else {
		if err != nil {
			return fmt.Sprintf("\033[31mError checking session: %v\033[0m", err)
		} else {
			return fmt.Sprintf("\033[33mUnknown status: %s\033[0m", raw)
		}
	}
}

// PrintOCIConfiguration displays the current configuration details
// and checks if required environment variables are set
func PrintOCIConfiguration() {
	displayBanner()

	// Get session status and display it with configuration details
	sessionStatus := CheckOCISessionValidity()
	fmt.Printf("\033[1mConfiguration Details:\033[0m %s\n", sessionStatus)

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
	fmt.Printf("\t      \033[33mVersion\033[0m: \033[32m%s\033[0m\n", buildinfo.Version)
	fmt.Println()
}
