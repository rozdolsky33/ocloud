package display

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/rozdolsky33/ocloud/buildinfo"
	"github.com/rozdolsky33/ocloud/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

var (
	boldStyle    = color.New(color.Bold)
	redStyle     = color.New(color.FgRed)
	greenStyle   = color.New(color.FgGreen)
	yellowStyle  = color.New(color.FgYellow)
	regularStyle = color.New(color.FgWhite)
)

var validRe = regexp.MustCompile(`(?i)^Session is valid until\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s*$`)
var expiredRe = regexp.MustCompile(`(?i)^Session has expired\s*$`)

// CheckOCISessionValidity checks the validity of the OCI session
// It assumes profile is not empty and should be called only with a valid profile
func CheckOCISessionValidity(profile string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "oci", "session", "validate", "--profile", profile)
	out, err := cmd.CombinedOutput()
	raw := strings.TrimSpace(string(out))

	if matches := validRe.FindStringSubmatch(raw); len(matches) > 1 {
		return greenStyle.Sprintf("Valid until %s", matches[1])
	} else if expiredRe.MatchString(raw) {
		return redStyle.Sprint("Session Expired")
	} else {
		if err != nil {
			return redStyle.Sprintf("Error checking session: %v", err)
		} else {
			return yellowStyle.Sprintf("Unknown status: %s", raw)
		}
	}
}

// RefresherStatus represents the status of the OCI auth refresher
type RefresherStatus struct {
	IsRunning bool
	PID       string
	Display   string
}

// CheckOCIAuthRefresherStatus checks if the OCI auth refresher script is running for the current profile
func CheckOCIAuthRefresherStatus() RefresherStatus {
	profile := os.Getenv(flags.EnvKeyProfile)
	if profile == "" {
		return RefresherStatus{
			IsRunning: false,
			PID:       "",
			Display:   redStyle.Sprint("OFF"),
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return RefresherStatus{
			IsRunning: false,
			PID:       "",
			Display:   redStyle.Sprint("OFF"),
		}
	}

	pidFilePath := filepath.Join(homeDir, flags.OCIConfigDirName, "sessions", profile, "refresher.pid")

	pidBytes, err := os.ReadFile(pidFilePath)
	if err != nil {
		return RefresherStatus{
			IsRunning: false,
			PID:       "",
			Display:   redStyle.Sprint("OFF"),
		}
	}

	pidStr := strings.TrimSpace(string(pidBytes))

	// Check if the process with this PID is running
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the process exists
	cmd := exec.CommandContext(ctx, "ps", "-p", pidStr, "-o", "pid=")
	if err := cmd.Run(); err != nil {
		_ = os.Remove(pidFilePath)
		return RefresherStatus{
			IsRunning: false,
			PID:       "",
			Display:   redStyle.Sprint("OFF"),
		}
	}

	// Then check if it's actually the refresher script for this profile
	cmd = exec.CommandContext(ctx, "pgrep", "-af", fmt.Sprintf("oci_auth_refresher.sh.*%s", profile))
	out, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(out))

	if err == nil && len(outStr) > 0 && strings.Contains(outStr, pidStr) {
		return RefresherStatus{
			IsRunning: true,
			PID:       pidStr,
			Display:   greenStyle.Sprintf("ON [%s]", pidStr),
		}
	}

	// Process exists, but it's not the refresher script for this profile
	// Remove the PID file as it's stale
	_ = os.Remove(pidFilePath)
	return RefresherStatus{
		IsRunning: false,
		PID:       "",
		Display:   redStyle.Sprint("OFF"),
	}
}

// PrintOCIConfiguration displays the current configuration details
func PrintOCIConfiguration() {
	displayBanner()

	profile := os.Getenv(flags.EnvKeyProfile)

	// Handle session status and profile display together to avoid redundancy
	var sessionStatus string
	if profile == "" {
		sessionStatus = redStyle.Sprint("Not set - Please set profile")
		fmt.Printf("%s %s\n", boldStyle.Sprint("Configuration Details:"), sessionStatus)
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyProfile), sessionStatus)
	} else {
		sessionStatus = CheckOCISessionValidity(profile)
		fmt.Printf("%s %s\n", boldStyle.Sprint("Configuration Details:"), sessionStatus)
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyProfile), profile)
	}

	tenancyName := os.Getenv(flags.EnvKeyTenancyName)
	if tenancyName == "" {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyName), redStyle.Sprint("Not set - Please set tenancy"))
	} else {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyName), tenancyName)
	}

	compartment := os.Getenv(flags.EnvKeyCompartment)
	if compartment == "" {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyCompartment), redStyle.Sprint("Not set - Please set compartmen name"))
	} else {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyCompartment), compartment)
	}

	refresherStatus := CheckOCIAuthRefresherStatus()
	fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyAutoRefresher), refresherStatus.Display)

	path := config.TenancyMapPath()
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyMapPath), redStyle.Sprint("Not set (file not found)"))
	} else if err != nil {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyMapPath), redStyle.Sprintf("Error checking file: %v", err))
	} else {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyMapPath), path)
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
	fmt.Printf("\t      %s: %s\n", regularStyle.Sprint("Version"), regularStyle.Sprint(buildinfo.Version))
	fmt.Println()
}
