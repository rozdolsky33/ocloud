package cmd

import (
	"fmt"
	"os"

	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debugMode bool
	rootCmd   = &cobra.Command{
		Use:               "ocloud",
		Short:             "Interact with Oracle Cloud Infrastructure",
		PersistentPreRunE: initializeConfig,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		SilenceErrors: true, // Don't print errors (we'll handle them)
		SilenceUsage:  true, // Don't print usage on error
	}
	// osExit is a variable to allow mocking in tests
	osExit = os.Exit
)

func initializeConfig(cmd *cobra.Command, args []string) error {
	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug logging enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	// TENANCY-ID: flag > ENV OCI_CLI_TENANCY > ENV OCI_TENANCY_NAME > OCI config file
	switch {
	case cmd.Flags().Changed(FlagNameTenancyID):
		tenancyID := viper.GetString(FlagNameTenancyID)
		logrus.Debugf("using tenancy OCID from falg %s: %s", FlagNameTenancyID, tenancyID)

	case os.Getenv(EnvOCITenancy) != "":
		tenancyID := os.Getenv(EnvOCITenancy)
		viper.Set(FlagNameTenancyID, tenancyID)
		logrus.Debugf("using tenancy OCID from env %s: %s", EnvOCITenancy, tenancyID)

	case os.Getenv(EnvOCITenancyName) != "":
		name := os.Getenv(EnvOCITenancyName)
		tenancyID, err := config.LookUpTenancyID(name)
		if err != nil {
			return fmt.Errorf("could not look up tenancy ID for %q: %w", name, err)
		}
		viper.Set(FlagNameTenancyID, tenancyID)
		logrus.Debugf("using tenancy OCID for env name %q: %s", name, tenancyID)

	default:
		if err := setUpTenancyFromOciConfigFile(); err != nil {
			return fmt.Errorf("could not load tenancy OCID: %w", err)
		}
	}

	// COMPARTMENT: flag > ENV OCI_COMPARTMENT
	switch {
	case cmd.Flags().Changed(FlagNameCompartment):
		comp := viper.GetString(FlagNameCompartment)
		logrus.Debugf("using compartment from flag: %s", comp)

	case os.Getenv(EnvOCICompartment) != "":
		comp := os.Getenv(EnvOCICompartment)
		viper.Set(FlagNameCompartment, comp)
		logrus.Debugf("using compartment from env %s: %s", EnvOCICompartment, comp)
	}

	return nil
}

func setUpTenancyFromOciConfigFile() error {
	tenancyID, err := config.GetTenancyOCID()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}
	logrus.Debugf("using tenancy OCID from OCI config file: %s", tenancyID)
	viper.SetDefault(FlagNameTenancyID, tenancyID)

	return nil
}

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
