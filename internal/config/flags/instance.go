// Package flags defines flag types and domain-specific flag collections for the CLI.
package flags

import "github.com/spf13/cobra"

var FlagDefaultLimit = 20
var FlagDefaultPage = 1

// Instance flags
var (
	ImageDetailsFlag = BoolFlag{
		Name:      FlagNameImageDetails,
		Shorthand: FlagShortImageDetails,
		Default:   false,
		Usage:     FlagDescImageDetails,
	}

	LimitFlag = IntFlag{
		Name:      FlagNameLimit,
		Shorthand: FlagShortLimit,
		Default:   FlagDefaultLimit,
		Usage:     FlagDescLimit,
	}

	PageFlag = IntFlag{
		Name:      FlagNamePage,
		Shorthand: FlagShortPage,
		Default:   FlagDefaultPage,
		Usage:     FlagDescPage,
	}

	JSONFlag = BoolFlag{
		Name:      FlagNameJSON,
		Shorthand: FlagShortJSON,
		Default:   false,
		Usage:     FlagDescJSON,
	}
)

// instanceFlags is a slice of all instance-related flags for batch registration
// This is kept for backward compatibility but is no longer used in the new subcommand structure
var instanceFlags = []Flag{
	ImageDetailsFlag,
	LimitFlag,
	PageFlag,
	JSONFlag,
}

// AddInstanceFlags adds all instance-related flags to the given command
// This is kept for backward compatibility but is no longer used in the new subcommand structure
func AddInstanceFlags(cmd *cobra.Command) {
	for _, f := range instanceFlags {
		f.Add(cmd)
	}
}
