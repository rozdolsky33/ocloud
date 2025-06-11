package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debugMode bool

	// rootCmd represents the base command used to interact with Oracle Cloud Infrastructure through the CLI.
	rootCmd = &cobra.Command{
		Use:               "ocloud",
		Short:             "Interact with Oracle Cloud Infrastructure",
		PersistentPreRunE: initConfig,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		SilenceErrors: true, // Don't print errors (we'll handle them)
		SilenceUsage:  true, // Don't print usage on error
	}
	// osExit is a variable to allow mocking in tests
	osExit = os.Exit
)

// init configures persistent flags and binds them to viper for managing application settings.
func init() {
	// debug flag
	rootCmd.PersistentFlags().
		BoolVarP(&debugMode, FlagNameDebug, FlagShortDebug, false, FlagDescDebug)

	// tenancy and compartment flags
	rootCmd.PersistentFlags().
		StringP(FlagNameTenancyID, FlagShortTenancyID, "", FlagDescTenancyID)
	rootCmd.PersistentFlags().
		StringP(FlagNameCompartment, FlagShortCompartment, "", FlagDescCompartment)

	// bind flags to viper keys
	_ = viper.BindPFlag(FlagNameTenancyID, rootCmd.PersistentFlags().Lookup(FlagNameTenancyID))
	_ = viper.BindPFlag(FlagNameCompartment, rootCmd.PersistentFlags().Lookup(FlagNameCompartment))

	// allow ENV overrides, e.g., OCI_CLI_TENANCY, OCI_TENANCY_NAME, OCI_COMPARTMENT
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()

}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		osExit(1)
	}
}
