package cmd

import (
	"os"
)

// isNoContextCommand provides functionality to check if a command doesn't need a full application context
func isNoContextCommand() bool {
	args := os.Args
	// If no arguments are provided, we don't need context
	if len(args) < 2 {
		return true
	}

	// Commands that don't need context
	noContextCommands := map[string]bool{
		"version": true,
		"config":  true,
	}

	// Flags that don't need context
	noContextFlags := map[string]bool{
		"--version": true,
		"-v":        true,
	}

	if noContextCommands[args[1]] {
		return true
	}

	for _, arg := range args[1:] {
		if noContextFlags[arg] {
			return true
		}
	}

	return false
}
