package flags

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var (
	RealmFlag = flags.StringFlag{
		Name:      flags.FlagNameRealm,
		Shorthand: flags.FlagShortRealm,
		Default:   "",
		Usage:     flags.FlagDescRealm,
	}

	// FilterFlag is used to filter regions by prefix
	FilterFlag = flags.StringFlag{
		Name:      flags.FlagNameFilter,
		Shorthand: flags.FlagShortFilter,
		Default:   "",
		Usage:     flags.FlagDescFilter,
	}
)
