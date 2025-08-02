package auth

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/rozdolsky33/ocloud/internal/logger"
)

//go:embed oci_auth_refresher.sh
var ociAuthRefresherScriptFS embed.FS

// ociAuthRefresherScript holds the content of the embedded script
var ociAuthRefresherScript string

// init initializes the ociAuthRefresherScript variable with the content of the embedded script file
func init() {
	scriptBytes, err := ociAuthRefresherScriptFS.ReadFile("oci_auth_refresher.sh")
	if err != nil {
		logger.LogWithLevel(logger.Logger, 0, "Failed to read embedded OCI auth refresher script", "error", err)
		return
	}
	ociAuthRefresherScript = string(scriptBytes)
}

// RunOCIAuthRefresher runs the OCI auth refresher script for the specified profile.
// The script keeps an OCI CLI session alive by refreshing it shortly before it expires.
func RunOCIAuthRefresher(profile string) error {
	logger.LogWithLevel(logger.Logger, 1, "Running OCI auth refresher script", "profile", profile)

	// Create a permanent directory for the script in the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	scriptDir := fmt.Sprintf("%s/.oci/scripts", homeDir)
	if err := os.MkdirAll(scriptDir, 0755); err != nil {
		return fmt.Errorf("failed to create script directory: %w", err)
	}

	scriptPath := fmt.Sprintf("%s/oci_auth_refresher.sh", scriptDir)

	// Write the script to the file
	if err := os.WriteFile(scriptPath, []byte(ociAuthRefresherScript), 0700); err != nil {
		return fmt.Errorf("failed to write OCI auth refresher script to file: %w", err)
	}

	// Use a background context to allow the script to run indefinitely
	ctx := context.Background()

	// Run the script using bash to ensure it's properly executed and visible to pgrep
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("NOHUP=1 %s %s", scriptPath, profile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start OCI auth refresher script: %w", err)
	}

	// Don't wait for the command to complete since it runs in the background
	pid := cmd.Process.Pid
	logger.LogWithLevel(logger.Logger, 1, "OCI auth refresher script started", "profile", profile, "pid", pid)

	// Print information about the process to the user
	fmt.Printf("\nOCI auth refresher started for profile %s with PID %d\n", profile, pid)
	fmt.Printf("You can verify it's running with: pgrep -af oci_auth_refresher.sh\n\n")

	return nil
}
