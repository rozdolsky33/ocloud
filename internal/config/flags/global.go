// Package flags defines flag types and domain-specific flag collections for the CLI.
package flags

import (
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global flags
var (
	LogLevelFlag = StringFlag{
		Name:    FlagNameLogLevel,
		Default: "info",
		Usage:   FlagDescLogLevel,
	}

	ColorFlag = BoolFlag{
		Name:    "color",
		Default: false,
		Usage:   logger.ColoredOutputMsg,
	}

	TenancyIDFlag = StringFlag{
		Name:      FlagNameTenancyID,
		Shorthand: FlagShortTenancyID,
		Default:   "",
		Usage:     FlagDescTenancyID,
	}

	TenancyNameFlag = StringFlag{
		Name:    FlagNameTenancyName,
		Default: "",
		Usage:   FlagDescTenancyName,
	}

	CompartmentFlag = StringFlag{
		Name:      FlagNameCompartment,
		Shorthand: FlagShortCompartment,
		Default:   "",
		Usage:     FlagDescCompartment,
	}
)

// globalFlags is a slice of all global flags for batch registration
var globalFlags = []Flag{
	LogLevelFlag,
	ColorFlag,
	TenancyIDFlag,
	TenancyNameFlag,
	CompartmentFlag,
}

// AddGlobalFlags adds all global flags to the given command
func AddGlobalFlags(cmd *cobra.Command) {
	// Add global flags as persistent flags
	cmd.PersistentFlags().StringP(LogLevelFlag.Name, LogLevelFlag.Shorthand, LogLevelFlag.Default, LogLevelFlag.Usage)
	cmd.PersistentFlags().BoolP(ColorFlag.Name, ColorFlag.Shorthand, ColorFlag.Default, ColorFlag.Usage)
	cmd.PersistentFlags().StringP(TenancyIDFlag.Name, TenancyIDFlag.Shorthand, TenancyIDFlag.Default, TenancyIDFlag.Usage)
	cmd.PersistentFlags().StringP(TenancyNameFlag.Name, TenancyNameFlag.Shorthand, TenancyNameFlag.Default, TenancyNameFlag.Usage)
	cmd.PersistentFlags().StringP(CompartmentFlag.Name, CompartmentFlag.Shorthand, CompartmentFlag.Default, CompartmentFlag.Usage)

	// Bind flags to viper for configuration
	_ = viper.BindPFlag(FlagNameTenancyID, cmd.PersistentFlags().Lookup(FlagNameTenancyID))
	_ = viper.BindPFlag(FlagNameTenancyName, cmd.PersistentFlags().Lookup(FlagNameTenancyName))
	_ = viper.BindPFlag(FlagNameCompartment, cmd.PersistentFlags().Lookup(FlagNameCompartment))

	// allow ENV overrides, e.g., OCI_CLI_TENANCY, OCI_TENANCY_NAME, OCI_COMPARTMENT
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
}
