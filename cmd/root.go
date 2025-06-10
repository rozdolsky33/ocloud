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
	}
)

// initializeConfig sets up logging and loads the tenancy OCID.
func initializeConfig(cmd *cobra.Command, args []string) error {
	// 1) log level
	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug logging enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// 2) fetch tenancy
	tenancyID, err := config.GetTenancyOCID()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}
	logrus.Debug("using tenancy OCID: ", tenancyID)

	// 3) make it available to flags/config
	viper.SetDefault(FlagNameTenancyID, tenancyID)
	return nil
}

func init() {
	// debug flag
	rootCmd.PersistentFlags().BoolVarP(&debugMode, FlagNameDebug, FlagShortDebug, false, FlagDescDebug)
	rootCmd.PersistentFlags().StringP(FlagNameTenancyID, FlagShortTenancyID, "", FlagDescTenancyID)
	rootCmd.PersistentFlags().StringP(FlagNameCompartment, FlagShortCompartment, "", FlagDescCompartment)

	_ = viper.BindPFlag(FlagNameTenancyID, rootCmd.PersistentFlags().Lookup(FlagNameTenancyID))
	_ = viper.BindPFlag(FlagNameCompartment, rootCmd.PersistentFlags().Lookup(FlagNameCompartment))
}

// Execute runs the CLI.
func Execute() {
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
