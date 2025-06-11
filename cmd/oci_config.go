package cmd

import (
	"fmt"
	"os"

	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initConfig configures logging levels and loads OCI tenancy or compartment details from flags, environment, or config files.
func initConfig(cmd *cobra.Command, args []string) error {
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
		tenancyID, err := config.LookupTenancyID(name)
		if err != nil {
			return fmt.Errorf("could not look up tenancy ID for %q: %w", name, err)
		}
		viper.Set(FlagNameTenancyID, tenancyID)
		logrus.Debugf("using tenancy OCID for env name %q: %s", name, tenancyID)

	default:
		if err := loadTenancyOCID(); err != nil {
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

// loadTenancyOCID loads the tenancy OCID from an OCI config file and sets it as the default value in viper.
// Returns an error if the tenancy OCID cannot be retrieved or there is an issue with the OCI config file.
func loadTenancyOCID() error {
	tenancyID, err := config.GetTenancyOCID()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}
	logrus.Debugf("using tenancy OCID from OCI config file: %s", tenancyID)
	viper.SetDefault(FlagNameTenancyID, tenancyID)

	return nil
}
