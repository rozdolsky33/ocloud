// Package bastion Shell/exec helpers. Keep all child-process spawning here so cancellation and I/O
// are consistent across the feature.
package bastion

import (
	"context"
	"io"
	"os/exec"
)

// RunShell runs the given command line using `bash -lc` and ties its lifetime to ctx.
// Stdout/Stderr are wired; Stdin is inherited from the process (interactive SSH).
func RunShell(ctx context.Context, stdout, stderr io.Writer, cmdLine string) error {
	cmd := exec.CommandContext(ctx, "bash", "-lc", cmdLine)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = nil // set to os.Stdin if you want interactive; flows can override if needed
	return cmd.Run()
}
