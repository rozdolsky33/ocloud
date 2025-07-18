package flags

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var (
	RealmFlag = flags.StringFlag{
		Name:      flags.FlagNameRealm,
		Shorthand: flags.FlagShortRealm,
		Default:   "",
		Usage:     flags.FlagDescRealm,
	}
)
