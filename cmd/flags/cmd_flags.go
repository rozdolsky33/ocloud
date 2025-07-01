package flags

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var FlagDefaultLimit = 20
var FlagDefaultPage = 1

var (
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

	SortFlag = flags.StringFlag{
		Name:      flags.FlagNameSort,
		Shorthand: flags.FlagShortSort,
		Default:   "",
		Usage:     flags.FlagDescSort,
	}
)
