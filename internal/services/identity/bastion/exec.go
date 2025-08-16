package bastion

// Package-level shell execution helper for bastion-related flows.
// Keeping child-process spawning in the service layer ensures CLI code remains
// thin and focused on user interaction while services encapsulate the execution
// details. This also centralizes context-aware process handling.

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
