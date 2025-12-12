package display

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rozdolsky33/ocloud/buildinfo"
	"github.com/rozdolsky33/ocloud/internal/config"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

var (
	boldStyle    = color.New(color.Bold)
	redStyle     = color.New(color.FgRed)
	greenStyle   = color.New(color.FgGreen)
	yellowStyle  = color.New(color.FgYellow)
	regularStyle = color.New(color.FgWhite)
)

// jwtClaims represents the claims in a JWT token
type jwtClaims struct {
	Exp int64 `json:"exp"`
}

// getSecurityTokenFile reads the OCI config file and extracts the security_token_file path for the given profile.
func getSecurityTokenFile(profile string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCIConfigFileName)
	file, err := os.Open(configPath)
	if err != nil {
		return "", fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inTargetProfile := false
	targetHeader := "[" + profile + "]"

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for the profile header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inTargetProfile = (line == targetHeader)
			continue
		}

		// If we're in the target profile, look for security_token_file
		if inTargetProfile {
			if strings.HasPrefix(line, "security_token_file") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan config file: %w", err)
	}

	// Fallback to a default path if not found in config
	return filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCISessionsDirName, profile, "token"), nil
}

// sessionExpiryFromTokenFile reads a JWT token file and extracts the expiration time.
func sessionExpiryFromTokenFile(tokenPath string) (time.Time, error) {
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("read token file: %w", err)
	}

	token := strings.TrimSpace(string(data))
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid token format")
	}

	payload := parts[1]
	// base64url without padding -> add padding if needed
	if m := len(payload) % 4; m != 0 {
		payload += strings.Repeat("=", 4-m)
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return time.Time{}, fmt.Errorf("decode payload: %w", err)
	}

	var claims jwtClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return time.Time{}, fmt.Errorf("unmarshal payload: %w", err)
	}
	if claims.Exp == 0 {
		return time.Time{}, fmt.Errorf("no exp claim in token")
	}

	return time.Unix(claims.Exp, 0), nil
}

// CheckOCISessionValidity checks the validity of the OCI session by parsing the JWT token directly.
// This avoids issues with debug output from the OCI CLI.
func CheckOCISessionValidity(profile string) string {
	tokenPath, err := getSecurityTokenFile(profile)
	if err != nil {
		return yellowStyle.Sprintf("Cannot find token: %v", err)
	}

	exp, err := sessionExpiryFromTokenFile(tokenPath)
	if err != nil {
		return yellowStyle.Sprintf("Cannot parse session token: %v", err)
	}

	now := time.Now()
	if now.After(exp) {
		return redStyle.Sprint("Session Expired")
	}

	// Format the expiry time in the same format as the old OCI CLI output
	ts := exp.Local().Format("2006-01-02 15:04:05")
	return greenStyle.Sprintf("Valid until %s", ts)
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

	pidFilePath := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCISessionsDirName, profile, flags.OCIRefresherPIDFile)

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

// PortForwardingStatus represents the status of active port-forwarding sessions
type PortForwardingStatus struct {
	IsActive bool
	Ports    []int
	Display  string
}

// CheckPortForwardingStatus checks for active SSH port-forwarding tunnels
func CheckPortForwardingStatus() PortForwardingStatus {
	tunnels, err := bastionSvc.GetActiveTunnels()
	if err != nil || len(tunnels) == 0 {
		return PortForwardingStatus{
			IsActive: false,
			Ports:    []int{},
			Display:  redStyle.Sprint("OFF"),
		}
	}

	var ports []int
	for _, tunnel := range tunnels {
		ports = append(ports, tunnel.LocalPort)
	}
	sort.Ints(ports)

	portStrs := make([]string, len(ports))
	for i, port := range ports {
		portStrs[i] = fmt.Sprintf("%d", port)
	}
	portsDisplay := strings.Join(portStrs, ", ")

	return PortForwardingStatus{
		IsActive: true,
		Ports:    ports,
		Display:  greenStyle.Sprintf("ON [%s]", portsDisplay),
	}
}

// PrintOCIConfiguration displays the current configuration details
func PrintOCIConfiguration() {
	displayBanner()

	profile := os.Getenv(flags.EnvKeyProfile)

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

	region, err := config.LoadOCIConfig().Region()

	if profile == "" {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyRegion), redStyle.Sprint("Not set - Please set profile first"))
	} else if err != nil {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyRegion), redStyle.Sprintf("Error loading region: %v", err))
	} else {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyRegion), region)
	}

	tenancyName := os.Getenv(flags.EnvKeyTenancyName)
	if tenancyName == "" {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyName), redStyle.Sprint("Not set - Please set tenancy"))
	} else {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyTenancyName), tenancyName)
	}

	compartment := os.Getenv(flags.EnvKeyCompartment)
	if compartment == "" {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyCompartment), redStyle.Sprint("Not set - Please set compartment name"))
	} else {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyCompartment), compartment)
	}

	refresherStatus := CheckOCIAuthRefresherStatus()
	fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyAutoRefresher), refresherStatus.Display)

	portForwardingStatus := CheckPortForwardingStatus()
	if portForwardingStatus.IsActive {
		fmt.Printf("  %s: %s\n", yellowStyle.Sprint(flags.EnvKeyPortForwarding), portForwardingStatus.Display)
	}

	path := config.TenancyMapPath()
	_, err = os.Stat(path)

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
