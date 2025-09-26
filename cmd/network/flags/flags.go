package flags

import "github.com/rozdolsky33/ocloud/internal/config/flags"

var (
	Gateway = flags.BoolFlag{
		Name:      flags.FlagNameGateway,
		Shorthand: flags.FlagShortGateway,
		Default:   false,
		Usage:     flags.FlagDescGateway,
	}
	Subnet = flags.BoolFlag{
		Name:      flags.FlagNameSubnet,
		Shorthand: flags.FlagShortSubnet,
		Default:   false,
		Usage:     flags.FlagDescSubnet,
	}
	Nsg = flags.BoolFlag{
		Name:      flags.FlagNameNsg,
		Shorthand: flags.FlagShortNsg,
		Default:   false,
		Usage:     flags.FlagDescNsg,
	}
	RouteTable = flags.BoolFlag{
		Name:      flags.FlagNameRoute,
		Shorthand: flags.FlagShortRoute,
		Default:   false,
		Usage:     flags.FlagDescRoute,
	}
	SecurityList = flags.BoolFlag{
		Name:      flags.FlagNameSecurity,
		Shorthand: flags.FlagShortSecurity,
		Default:   false,
		Usage:     flags.FlagDescSecurity,
	}
)
