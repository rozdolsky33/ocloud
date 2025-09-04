// Package flags define flag types and domain-specific flag collections for the CLI.
package flags

import "github.com/rozdolsky33/ocloud/internal/config/flags"

// AllInfoFlag Compute flags
var (
	AllInfoFlag = flags.BoolFlag{
		Name:      flags.FlagNameAllInformation,
		Shorthand: flags.FlagShortAllInformation,
		Default:   false,
		Usage:     flags.FlagDescAllInformation,
	}
)
