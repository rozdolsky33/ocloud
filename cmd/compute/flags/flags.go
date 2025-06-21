// Package flags defines flag types and domain-specific flag collections for the CLI.
package flags

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var FlagDefaultLimit = 20
var FlagDefaultPage = 1

// Compute flags
var (
	ImageDetailsFlag = flags.BoolFlag{
		Name:      flags.FlagNameImageDetails,
		Shorthand: flags.FlagShortImageDetails,
		Default:   false,
		Usage:     flags.FlagDescImageDetails,
	}

	LimitFlag = flags.IntFlag{
		Name:      flags.FlagNameLimit,
		Shorthand: flags.FlagShortLimit,
		Default:   FlagDefaultLimit,
		Usage:     flags.FlagDescLimit,
	}

	PageFlag = flags.IntFlag{
		Name:      flags.FlagNamePage,
		Shorthand: flags.FlagShortPage,
		Default:   FlagDefaultPage,
		Usage:     flags.FlagDescPage,
	}
)
