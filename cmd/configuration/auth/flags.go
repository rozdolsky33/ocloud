package auth

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var (
	// FilterFlag is used to filter regions by prefix
	FilterFlag = flags.StringFlag{
		Name:      flags.FlagNameFilter,
		Shorthand: flags.FlagShortFilter,
		Default:   "",
		Usage:     flags.FlagDescFilter,
	}
)
