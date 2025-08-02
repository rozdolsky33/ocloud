package auth

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/scripts"
)

// RunOCIAuthRefresher runs the OCI auth refresher script for the specified profile.
func RunOCIAuthRefresher(profile string) error {
	logger.LogWithLevel(logger.Logger, 1, "Running OCI auth refresher script", "profile", profile)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	scriptDir := fmt.Sprintf("%s/.oci/scripts", homeDir)
	if err := os.MkdirAll(scriptDir, 0o755); err != nil {
		return fmt.Errorf("failed to create script directory: %w", err)
	}

	scriptPath := fmt.Sprintf("%s/oci_auth_refresher.sh", scriptDir)

	// Write the embedded script bytes to the disk
	if err := os.WriteFile(scriptPath, scripts.OCIAuthRefresher, 0o700); err != nil {
		return fmt.Errorf("failed to write OCI auth refresher script to file: %w", err)
	}

	// Use a background context so it can run indefinitely
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("NOHUP=1 %s %s", scriptPath, profile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start OCI auth refresher script: %w", err)
	}

	pid := cmd.Process.Pid
	logger.LogWithLevel(logger.Logger, 1, "OCI auth refresher script started", "profile", profile, "pid", pid)

	fmt.Printf("\nOCI auth refresher started for profile %s with PID %d\n", profile, pid)
	fmt.Printf("You can verify it's running with: pgrep -af oci_auth_refresher.sh\n\n")
	return nil
}
