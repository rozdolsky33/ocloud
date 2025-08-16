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

//func SpawnDetached(cmdLine string) error {
//
//	cmd := exec.Command("ssh", cmdLine)
//	// Detach from this terminal/session so it survives after we exit.
//	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
//	// Don't inherit our stdio; log to a file instead.
//	f, err := os.OpenFile("/tmp/ssh-tunnel.log",
//		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//	cmd.Stdout = f
//	cmd.Stderr = f
//	cmd.Stdin = nil
//
//	// Start and release so we don't wait on it.
//	if err := cmd.Start(); err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("spawned tunnel pid=%d", cmd.Process.Pid)
//	_ = cmd.Process.Release()
//
//	return nil
//}
