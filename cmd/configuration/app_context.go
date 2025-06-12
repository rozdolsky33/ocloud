package configuration

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// AppContext represents the application context containing OCI configuration, tenancy, and compartment information.
type AppContext struct {
	Provider       common.ConfigurationProvider
	IdentityClient identity.IdentityClient
	TenancyID      string
	Compartment    string
	CompartmentID  string
}

// InitGlobalFlags initializes global CLI flags and binds them to environment variables for configuration.
func InitGlobalFlags(root *cobra.Command) {
	root.PersistentFlags().StringVarP(&helpers.LogLevel, FlagNameLogLevel, "", "info", helpers.LogLevelMsg)
	root.PersistentFlags().BoolVar(&helpers.ColoredOutput, "color", false, helpers.ColoredOutputMsg)
	root.PersistentFlags().StringP(FlagNameTenancyID, FlagShortTenancyID, "", FlagDescTenancyID)
	root.PersistentFlags().StringP(FlagNameCompartment, FlagShortCompartment, "", FlagDescCompartment)

	_ = viper.BindPFlag(FlagNameTenancyID, root.PersistentFlags().Lookup(FlagNameTenancyID))
	_ = viper.BindPFlag(FlagNameCompartment, root.PersistentFlags().Lookup(FlagNameCompartment))

	// allow ENV overrides, e.g., OCI_CLI_TENANCY, OCI_TENANCY_NAME, OCI_COMPARTMENT
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
}

// NewAppContext builds an AppContext by resolving tenancy and compartment from flags, env vars, or OCI config,
// then initializes the necessary OCI clients.
func NewAppContext(cmd *cobra.Command, _ []string) (*AppContext, error) {
	logger := helpers.CmdLogger
	logger.Info("Initializing application context.....")

	// Load OCI shared config and create a provider
	prov := config.LoadOCIConfig()

	// TENANCY-ID: flag > ENV OCI_CLI_TENANCY > ENV OCI_TENANCY_NAME > Resolve via a config file
	var tenancyID string

	if cmd.Flags().Changed(FlagNameTenancyID) {
		tenancyID = viper.GetString(FlagNameTenancyID)
		logger.V(1).Info("using tenancy OCID from flag",
			"flag", FlagNameTenancyID, "tenancyID", tenancyID)
	} else if env := os.Getenv(EnvOCITenancy); env != "" {
		tenancyID = env
		viper.Set(FlagNameTenancyID, tenancyID)
		logger.V(1).Info("using tenancy OCID from env",
			"env", EnvOCITenancy, "tenancyID", tenancyID)
	} else if name := os.Getenv(EnvOCITenancyName); name != "" {
		id, err := config.LookupTenancyID(name)
		if err != nil {
			return nil, fmt.Errorf("could not look up tenancy ID for %q: %w", name, err)
		}
		tenancyID = id
		viper.Set(FlagNameTenancyID, tenancyID)
		logger.V(1).Info("using tenancy OCID for env name",
			"name", name, "tenancyID", tenancyID)
	} else {
		// fall back to ResolveTenancyID, which will load from an OCI config file if needed
		if err := loadTenancyOCID(); err != nil {
			return nil, fmt.Errorf("could not load tenancy OCID: %w", err)
		}
		tenancyID = viper.GetString(FlagNameTenancyID)
	}

	// COMPARTMENT: flag > ENV OCI_COMPARTMENT > viper default
	var compartment string

	if cmd.Flags().Changed(FlagNameCompartment) {
		compartment = viper.GetString(FlagNameCompartment)
		logger.V(1).Info("using compartment from flag",
			"compartment", compartment)
	} else if env := os.Getenv(EnvOCICompartment); env != "" {
		compartment = env
		viper.Set(FlagNameCompartment, compartment)
		logger.V(1).Info("using compartment from env",
			"env", EnvOCICompartment, "compartment", compartment)
	} else {
		compartment = viper.GetString(FlagNameCompartment)
		logger.V(1).Info("using compartment from default",
			"compartment", compartment)
	}

	// Initialize OCI service clients
	idc, err := identity.NewIdentityClientWithConfigurationProvider(prov)
	if err != nil {
		return nil, fmt.Errorf("creating identity client: %w", err)
	}

	// Optionally override region from env
	if region, ok := os.LookupEnv(EnvOCIRegion); ok {
		idc.SetRegion(region)
		logger.V(1).Info("overriding region from env", "region", region)
	}

	return &AppContext{
		Provider:       prov,
		IdentityClient: idc,
		TenancyID:      tenancyID,
		Compartment:    compartment,
	}, nil
}

// loadTenancyOCID loads the tenancy OCID from an OCI config file and sets it as the default value in viper.
// Returns an error if the tenancy OCID cannot be retrieved or there is an issue with the OCI config file.
func loadTenancyOCID() error {
	tenancyID, err := config.GetTenancyOCID()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}
	logger := helpers.CmdLogger
	logger.V(1).Info("using tenancy OCID from OCI config file", "tenancyID", tenancyID)
	viper.SetDefault(FlagNameTenancyID, tenancyID)

	return nil
}
