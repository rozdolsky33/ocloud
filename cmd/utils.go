package cmd

import (
	"os"
)

// isNoContextCommand provides functionality to check if a command doesn't need a full application context
// This is a simplified version of the previous CommandRegistry, removing unused methods
func isNoContextCommand() bool {
	args := os.Args
	// If no arguments are provided (just the program name), we don't need context
	// This avoids initialization when just displaying help/usage information
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

// isRootCommandWithoutSubcommands checks if the command being executed is the root command without any subcommands or flags
// This is used to determine whether to display the banner and configuration details
func isRootCommandWithoutSubcommands() bool {
	args := os.Args
	return len(args) == 1
}
