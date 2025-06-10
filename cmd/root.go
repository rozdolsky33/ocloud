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

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug logging enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	//--------------------------------tenancy OCID from OCI config file-----------------------------------------------//

	tenancyErr := setUpTenancyFromOciConfigFile()
	if tenancyErr != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", tenancyErr)
	}

	//--------------------------------tenancy OCID from OCI_CLI_TENANCY env-------------------------------------------//

	// Overwrite tenancy OCID with the value of OCI_CLI_TENANCY, if set.
	_, err := bindTenancyIDFromEnv()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}

	//-------------------------------tenancy OCID from OCI_TENANCY_NAME env-------------------------------------------//

	// Overwrite tenancy OCID with the value of OCI_CLI_TENANCY, if set.
	_, err = setTenancyIDFromEnvName()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}

	//-------------------------------tenancy OCID from OCI_COMPARTMENT env---------------------------------------------//

	// Overwrite the compartment name with the value of OCI_COMPARTMENT, if set.
	_, err = setCompartmentNameFromEnv()
	if err != nil {
		return fmt.Errorf("could not set comaprtmetn name from ENV: %w", err)
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

// bindTenancyFromEnv checks OCI_CLI_TENANCY and, if present,
// injects its value into viper under "tenancy-id".
// Returns true if we found & set the env-var, false otherwise.
func bindTenancyIDFromEnv() (bool, error) {
	tenancyID, exists := os.LookupEnv(EnvOCITenancy)
	if !exists {
		return false, nil
	}
	// set the actual OCID value into viper (the highest precedence)
	viper.Set(FlagNameTenancyID, tenancyID)
	logrus.Debugf("setting tenancy OCID from env %s: %s", EnvOCITenancy, tenancyID)
	logrus.Debugf("overwritten tenancy OCID from env: %v", exists)
	return true, nil
}

// setTenancyIDFromEnvName sets the tenancy ID based on the `OCI_TENANCY_NAME` environment variable if it's defined.
func setTenancyIDFromEnvName() (bool, error) {
	tenancyName, exists := os.LookupEnv(EnvOCITenancyName)
	if !exists {
		return false, nil
	}
	// If OCI_TENANCY_NAME is set, look up the tenancy ID from the tenancy map
	tenancyID, err := config.LookUpTenancyID(tenancyName)
	if err != nil {
		return false, fmt.Errorf("could not look up tenancy ID: %w", err)
	}
	logrus.Debug("using tenancy OCID: ", tenancyID)
	viper.Set(FlagNameTenancyID, tenancyID)
	return true, nil
}

// If OCI_COMPARTMENT is set, set the compartment name
func setCompartmentNameFromEnv() (bool, error) {
	compartmentName, exists := os.LookupEnv(EnvOCICompartment)
	if !exists {
		return false, nil
	}

	err := viper.BindEnv("compartment", EnvOCICompartment)
	if err != nil {
		return false, fmt.Errorf("could not bind env: %w", err)
	}
	logrus.Debugf("setting compartment from env %s: %s", EnvOCICompartment, compartmentName)
	logrus.Debugf("overwritten compartment from env: %v", exists)
	return true, nil
}

// init configures persistent flags and binds them to viper for managing application settings.
func init() {
	// debug flag
	rootCmd.PersistentFlags().BoolVarP(&debugMode, FlagNameDebug, FlagShortDebug, false, FlagDescDebug)
	var flagTenancyId string
	rootCmd.PersistentFlags().StringVarP(&flagTenancyId, FlagNameTenancyID, FlagShortTenancyID, "", FlagDescTenancyID)

	var flagCompartmentName string
	rootCmd.PersistentFlags().StringVarP(&flagCompartmentName, FlagNameCompartment, FlagShortCompartment, "", FlagDescCompartment)

	//_ = viper.BindPFlag(FlagNameTenancyID, rootCmd.PersistentFlags().Lookup(FlagNameTenancyID))
	//_ = viper.BindPFlag(FlagNameCompartment, rootCmd.PersistentFlags().Lookup(FlagNameCompartment))
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
