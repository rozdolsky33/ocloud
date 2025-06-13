// Package flags defines flag types and domain-specific flag collections for the CLI.
package flags

import (
	"github.com/spf13/cobra"
)

// BoolFlag represents a boolean command flag configuration
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

// StringFlag represents a string command flag configuration
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

// IntFlag represents an integer command flag configuration
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

// Flag is an interface that all flag types must implement
type Flag interface {
	Add(*cobra.Command)
}
