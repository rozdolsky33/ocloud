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

const (
	reset = "\033[0m"
	bold  = "\033[1m"

	red   = "\033[31m"
	green = "\033[32m"
	yel   = "\033[33m"
)

func colorize(s, color string) string { return color + s + reset }

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
		return colorize(fmt.Sprintf("Valid until %s", matches[1]), green)
	} else if expiredRe.MatchString(raw) {
		return colorize("Session Expired", red)
	} else {
		if err != nil {
			return colorize(fmt.Sprintf("Error checking session: %v", err), red)
		} else {
			return colorize(fmt.Sprintf("Unknown status: %s", raw), yel)
		}
	}
}

// PrintOCIConfiguration displays the current configuration details
// and checks if required environment variables are set
func PrintOCIConfiguration() {
	displayBanner()

	sessionStatus := CheckOCISessionValidity()
	fmt.Printf("%s %s\n", colorize("Configuration Details:", bold), sessionStatus)

	profile := os.Getenv("OCI_CLI_PROFILE")
	if profile == "" {
		fmt.Printf("  %s: %s\n", colorize("OCI_CLI_PROFILE", yel), colorize("Not set - Please set profile", red))
	} else {
		fmt.Printf("  %s: %s\n", colorize("OCI_CLI_PROFILE", yel), profile)
	}

	tenancyName := os.Getenv(flags.EnvOCITenancyName)
	if tenancyName == "" {
		fmt.Printf("  %s: %s\n", colorize("OCI_TENANCY_NAME", yel), colorize("Not set - Please set tenancy", red))
	} else {
		fmt.Printf("  %s: %s\n", colorize("OCI_TENANCY_NAME", yel), tenancyName)
	}

	compartment := os.Getenv(flags.EnvOCICompartment)
	if compartment == "" {
		fmt.Printf("  %s: %s\n", colorize("OCI_COMPARTMENT", yel), colorize("Not set - Please set compartmen name", red))
	} else {
		fmt.Printf("  %s: %s\n", colorize("OCI_COMPARTMENT", yel), compartment)
	}

	path := config.TenancyMapPath()
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Printf("  %s: %s\n", colorize("OCI_TENANCY_MAP_PATH", yel), colorize("Not set (file not found)", red))
	} else if err != nil {
		fmt.Printf("  %s: %s\n", colorize("OCI_TENANCY_MAP_PATH", yel), colorize(fmt.Sprintf("Error checking file: %v", err), red))
	} else {
		fmt.Printf("  %s: %s\n", colorize("OCI_TENANCY_MAP_PATH", yel), path)
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
	fmt.Printf("\t      %s: %s\n", colorize("Version", bold), colorize(buildinfo.Version, green))
	fmt.Println()
}
