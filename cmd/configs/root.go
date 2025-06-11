package configs

import (
	"github.com/rozdolsky33/ocloud/internal/helpers"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// ConfigCmd represents the base command used to interact with Oracle Cloud Infrastructure through the CLI.
	ConfigCmd = &cobra.Command{
		Use:          "config",
		Short:        "Configure Interact with Oracle Cloud Infrastructure",
		RunE:         initConfig,
		PreRunE:      preConfigE,
		SilenceUsage: true, // Don't print usage on error
	}
	// osExit is a variable to allow mocking in tests
	osExit = os.Exit
)

func preConfigE(cmd *cobra.Command, args []string) error {
	if err := helpers.SetLogger(); err != nil {
		return err
	}
	// Initialize the internal/config logger with the CmdLogger
	helpers.InitLogger(helpers.CmdLogger)

	return nil
}

// init configures persistent flags and binds them to viper for managing application settings.
func init() {
	// tenancy and compartment flags
	ConfigCmd.PersistentFlags().StringP(FlagNameTenancyID, FlagShortTenancyID, "", FlagDescTenancyID)
	ConfigCmd.PersistentFlags().StringP(FlagNameCompartment, FlagShortCompartment, "", FlagDescCompartment)

	// bind flags to viper keys
	_ = viper.BindPFlag(FlagNameTenancyID, ConfigCmd.PersistentFlags().Lookup(FlagNameTenancyID))
	_ = viper.BindPFlag(FlagNameCompartment, ConfigCmd.PersistentFlags().Lookup(FlagNameCompartment))

}
