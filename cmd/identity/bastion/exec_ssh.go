// Package bastion Shell/exec helpers. Keep all child-process spawning here so cancellation and I/O
// are consistent across the feature.
package bastion

import (
	"context"
	"io"
	"os"
	"os/exec"
)

// RunShell runs the given command line using `bash -lc` and ties its lifetime to ctx.
// Stdout/Stderr are wired; Stdin is inherited from the current process (enables interactive SSH by default).
func RunShell(ctx context.Context, stdout, stderr io.Writer, cmdLine string) error {
	cmd := exec.CommandContext(ctx, "bash", "-lc", cmdLine)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
