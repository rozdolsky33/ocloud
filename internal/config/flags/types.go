// Package flags provides a type-safe and reusable way to define and manage command-line flags
// for CLI applications using cobra and pflag libraries. It offers structured flag types for
// boolean, string, and integer values, along with consistent interfaces for adding these flags
// to commands and flag sets.
package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// BoolFlag represents a boolean command flag configuration with a name, optional shorthand,
// default value, and usage description. It implements the Flag interface for boolean flags.
type BoolFlag struct {
	Name      string
	Shorthand string
	Default   bool
	Usage     string
}

// Add adds the boolean flag to the command
func (f BoolFlag) Add(cmd *cobra.Command) {
	cmd.Flags().BoolP(f.Name, f.Shorthand, f.Default, f.Usage)
}

// Apply adds the boolean flag to the given flag set
func (f BoolFlag) Apply(flags *pflag.FlagSet) {
	flags.BoolP(f.Name, f.Shorthand, f.Default, f.Usage)
}

// StringFlag represents a string command flag configuration with a name, optional shorthand,
// default value, and usage description. It implements the Flag interface for string flags.
type StringFlag struct {
	Name      string
	Shorthand string
	Default   string
	Usage     string
}

// Add adds the string flag to the command
func (f StringFlag) Add(cmd *cobra.Command) {
	cmd.Flags().StringP(f.Name, f.Shorthand, f.Default, f.Usage)
}

// Apply adds the string flag to the given flag set
func (f StringFlag) Apply(flags *pflag.FlagSet) {
	flags.StringP(f.Name, f.Shorthand, f.Default, f.Usage)
}

// IntFlag represents an integer command flag configuration with a name, optional shorthand,
// default value, and usage description. It implements the Flag interface for integer flags.
type IntFlag struct {
	Name      string
	Shorthand string
	Default   int
	Usage     string
}

// Add adds the integer flag to the command
func (f IntFlag) Add(cmd *cobra.Command) {
	cmd.Flags().IntP(f.Name, f.Shorthand, f.Default, f.Usage)
}

// Apply adds the integer flag to the given flag set
func (f IntFlag) Apply(flags *pflag.FlagSet) {
	flags.IntP(f.Name, f.Shorthand, f.Default, f.Usage)
}

// Flag defines the interface that all flag types must implement to be used within the CLI.
// It provides methods for adding flags to both cobra.Command and pflag.FlagSet, allowing
// flexible flag registration across different command contexts.
type Flag interface {
	// Add registers the flag with the provided cobra.Command
	Add(*cobra.Command)
	// Apply registers the flag with the provided pflag.FlagSet
	Apply(*pflag.FlagSet)
}
