package cmdutil

import (
	"os"
)

// IsNoContextCommand checks if a command doesn't need a full application context
func IsNoContextCommand() bool {
	args := os.Args
	// If no arguments are provided (just the program name), we don't need context
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

// IsRootCommandWithoutSubcommands checks if the command being executed is the root command without any subcommands or flags
func IsRootCommandWithoutSubcommands() bool {
	args := os.Args
	return len(args) == 1
}
