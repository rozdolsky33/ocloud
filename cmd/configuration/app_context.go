package configuration

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/cmd/configs"
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
	root.PersistentFlags().StringVarP(&helpers.LogLevel, configs.FlagNameLogLevel, "", "info", helpers.LogLevelMsg)
	root.PersistentFlags().BoolVar(&helpers.ColoredOutput, "color", false, helpers.ColoredOutputMsg)
	root.PersistentFlags().StringP(configs.FlagNameTenancyID, configs.FlagShortTenancyID, "", configs.FlagDescTenancyID)
	root.PersistentFlags().StringP(configs.FlagNameCompartment, configs.FlagShortCompartment, "", configs.FlagDescCompartment)

	_ = viper.BindPFlag(configs.FlagNameTenancyID, root.PersistentFlags().Lookup(configs.FlagNameTenancyID))
	_ = viper.BindPFlag(configs.FlagNameCompartment, root.PersistentFlags().Lookup(configs.FlagNameCompartment))

	// allow ENV overrides, e.g., OCI_CLI_TENANCY, OCI_TENANCY_NAME, OCI_COMPARTMENT
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
}

// NewAppContext builds an AppContext by resolving tenancy and compartment from flags, env vars, or OCI config,
// then initializes the necessary OCI clients.
func NewAppContext(cmd *cobra.Command, _ []string) (*AppContext, error) {
	logger := helpers.CmdLogger
	logger.Info("App context initialization started.....")

	// Load OCI shared config and create a provider
	prov := config.LoadOCIConfig()

	// TENANCY-ID: flag > ENV OCI_CLI_TENANCY > ENV OCI_TENANCY_NAME > Resolve via a config file
	var tenancyID string

	if cmd.Flags().Changed(configs.FlagNameTenancyID) {
		tenancyID = viper.GetString(configs.FlagNameTenancyID)
		logger.V(1).Info("using tenancy OCID from flag",
			"flag", configs.FlagNameTenancyID, "tenancyID", tenancyID)
	} else if env := os.Getenv(configs.EnvOCITenancy); env != "" {
		tenancyID = env
		viper.Set(configs.FlagNameTenancyID, tenancyID)
		logger.V(1).Info("using tenancy OCID from env",
			"env", configs.EnvOCITenancy, "tenancyID", tenancyID)
	} else if name := os.Getenv(configs.EnvOCITenancyName); name != "" {
		id, err := config.LookupTenancyID(name)
		if err != nil {
			return nil, fmt.Errorf("could not look up tenancy ID for %q: %w", name, err)
		}
		tenancyID = id
		viper.Set(configs.FlagNameTenancyID, tenancyID)
		logger.V(1).Info("using tenancy OCID for env name",
			"name", name, "tenancyID", tenancyID)
	} else {
		// fall back to ResolveTenancyID, which will load from an OCI config file if needed
		if err := loadTenancyOCID(); err != nil {
			return nil, fmt.Errorf("could not load tenancy OCID: %w", err)
		}
		tenancyID = viper.GetString(configs.FlagNameTenancyID)
	}

	// COMPARTMENT: flag > ENV OCI_COMPARTMENT > viper default
	var compartment string

	if cmd.Flags().Changed(configs.FlagNameCompartment) {
		compartment = viper.GetString(configs.FlagNameCompartment)
		logger.V(1).Info("using compartment from flag",
			"compartment", compartment)
	} else if env := os.Getenv(configs.EnvOCICompartment); env != "" {
		compartment = env
		viper.Set(configs.FlagNameCompartment, compartment)
		logger.V(1).Info("using compartment from env",
			"env", configs.EnvOCICompartment, "compartment", compartment)
	} else {
		compartment = viper.GetString(configs.FlagNameCompartment)
		logger.V(1).Info("using compartment from default",
			"compartment", compartment)
	}

	// Initialize OCI service clients
	idc, err := identity.NewIdentityClientWithConfigurationProvider(prov)
	if err != nil {
		return nil, fmt.Errorf("creating identity client: %w", err)
	}

	// Optionally override region from env
	if region, ok := os.LookupEnv(configs.EnvOCIRegion); ok {
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
	viper.SetDefault(configs.FlagNameTenancyID, tenancyID)

	return nil
}
