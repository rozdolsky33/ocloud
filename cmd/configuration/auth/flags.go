package auth

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var (
	// EnvOnlyFlag is used to only output environment variables without running interactive authentication
	EnvOnlyFlag = flags.BoolFlag{
		Name:      flags.FlagNameEnvOnly,
		Shorthand: flags.FlagShortEnvOnly,
		Default:   false,
		Usage:     flags.FlagDescEnvOnly,
	}

	// FilterFlag is used to filter regions by prefix
	FilterFlag = flags.StringFlag{
		Name:      flags.FlagNameFilter,
		Shorthand: flags.FlagShortFilter,
		Default:   "",
		Usage:     flags.FlagDescFilter,
	}
)
