// Package flags defines flag types and domain-specific flag collections for the CLI.
// It provides utility functions for safely retrieving flag values from cobra commands
// with default fallbacks. This package is used throughout the application to handle
// command-line flags in a consistent manner.
//
// Usage:
//   - Import this package in your command handlers
//   - Use the Get*Flag functions to safely retrieve flag values
//   - Each function handles potential errors and returns a default value if the flag is not found
//
// Example:
//   debug := flags.GetBoolFlag(cmd, "debug", false)
//   name: = flags.GetStringFlag(cmd, "name", "default-name")

package flags

import (
	"github.com/spf13/cobra"
)

// GetBoolFlag retrieves a boolean flag value from the cobra.Command.
// It provides a safe way to access boolean flags with automatic error handling.
// If the flag is not found or there's an error reading it, it returns the provided default value.
//
// Parameters:
//   - cmd: The cobra.Command instance containing the flags
//   - flagName: The name of the flag to retrieve
//   - defaultValue: The value to return if the flag is not found or has an error
//
// Returns:
//   - The boolean value of the flag or the default value if not found/error
func GetBoolFlag(cmd *cobra.Command, flagName string, defaultValue bool) bool {
	value, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetStringFlag gets a string flag value from the command
// If the flag is not found or there's an error, it returns the default value
func GetStringFlag(cmd *cobra.Command, flagName string, defaultValue string) string {
	value, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetIntFlag gets an integer flag value from the command
// If the flag is not found or there's an error, it returns the default value
func GetIntFlag(cmd *cobra.Command, flagName string, defaultValue int) int {
	value, err := cmd.Flags().GetInt(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetStringSliceFlag gets a string slice flag value from the command
// If the flag is not found or there's an error, it returns the default value
func GetStringSliceFlag(cmd *cobra.Command, flagName string, defaultValue []string) []string {
	value, err := cmd.Flags().GetStringSlice(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetFloat64Flag gets a float64 flag value from the command
// If the flag is not found or there's an error, it returns the default value
func GetFloat64Flag(cmd *cobra.Command, flagName string, defaultValue float64) float64 {
	value, err := cmd.Flags().GetFloat64(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDurationFlag gets a duration flag value from the command
// If the flag is not found or there's an error, it returns the default value
func GetDurationFlag(cmd *cobra.Command, flagName string, defaultValue int64) int64 {
	value, err := cmd.Flags().GetInt64(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}
