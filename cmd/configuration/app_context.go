package configuration

import (
	"context"
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
// AppContext holds the application-wide OCI configuration and resolved IDs.
type AppContext struct {
	Ctx             context.Context
	Provider        common.ConfigurationProvider
	IdentityClient  identity.IdentityClient
	TenancyID       string
	TenancyName     string
	CompartmentName string
	CompartmentID   string
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

// NewAppContext initializes AppContext, resolves tenancy & compartment IDs, and builds OCI clients.
func NewAppContext(ctx context.Context, cmd *cobra.Command, _ []string) (*AppContext, error) {
	logger := helpers.CmdLogger
	logger.Info("Initializing application context")

	// Load OCI config & create provider
	prov := config.LoadOCIConfig()

	// Create an identity client (needed for compartment lookup)
	idClient, err := identity.NewIdentityClientWithConfigurationProvider(prov)
	if err != nil {
		return nil, fmt.Errorf("creating identity client: %w", err)
	}
	// Optional region override
	if region, ok := os.LookupEnv(EnvOCIRegion); ok {
		idClient.SetRegion(region)
		logger.V(1).Info("overriding region from env", "region", region)
	}

	// Build base AppContext
	appCtx := &AppContext{
		Ctx:             ctx,
		Provider:        prov,
		IdentityClient:  idClient,
		TenancyName:     viper.GetString(FlagNameTenancyName),
		CompartmentName: viper.GetString(FlagNameCompartment),
	}

	// Resolve Tenancy OCID: flag > ENV OCI_CLI_TENANCY > ENV OCI_TENANCY_NAME > OCI config file
	var tenancyID string
	envTenancy := os.Getenv(EnvOCITenancy)
	envTenancyName := os.Getenv(EnvOCITenancyName)

	switch {
	case cmd.Flags().Changed(FlagNameTenancyID):
		tenancyID = viper.GetString(FlagNameTenancyID)
		logger.V(1).Info("using tenancy OCID from flag", "tenancyID", tenancyID)

	case envTenancy != "":
		tenancyID = envTenancy
		viper.Set(FlagNameTenancyID, tenancyID)
		logger.V(1).Info("using tenancy OCID from env", "tenancyID", tenancyID)

	case envTenancyName != "":
		lookupID, err := config.LookupTenancyID(envTenancyName)
		if err != nil {
			return nil, fmt.Errorf("could not look up tenancy ID for %q: %w", envTenancyName, err)
		}
		tenancyID = lookupID
		viper.Set(FlagNameTenancyID, tenancyID)
		logger.V(1).Info("using tenancy OCID for name", "tenancyName", envTenancyName, "tenancyID", tenancyID)

	default:
		// load from OCI config file
		fileID, err := config.GetTenancyOCID()
		if err != nil {
			return nil, fmt.Errorf("could not load tenancy OCID: %w", err)
		}
		tenancyID = fileID
		logger.V(1).Info("using tenancy OCID from config file", "tenancyID", tenancyID)
	}
	viper.Set(FlagNameTenancyID, tenancyID)
	appCtx.TenancyID = tenancyID

	// Resolve Compartment OCID using helper
	compID, err := fetchCompartmentID(appCtx.Ctx, appCtx.TenancyID, appCtx.CompartmentName, appCtx.IdentityClient)
	if err != nil {
		return nil, fmt.Errorf("could not resolve compartment ID: %w", err)
	}
	appCtx.CompartmentID = compID

	return appCtx, nil
}

// fetchTenancyOCID loads the tenancy OCID from an OCI config file and sets it as the default value in viper.
// Returns an error if the tenancy OCID cannot be retrieved or there is an issue with the OCI config file.
func fetchTenancyOCID() error {
	tenancyID, err := config.GetTenancyOCID()
	if err != nil {
		return fmt.Errorf("could not load tenancy OCID: %w", err)
	}
	logger := helpers.CmdLogger
	logger.V(1).Info("using tenancy OCID from OCI config file", "tenancyID", tenancyID)
	viper.SetDefault(FlagNameTenancyID, tenancyID)

	return nil
}

func fetchCompartmentID(ctx context.Context, tenancyOCID, compartmentName string, idClient identity.IdentityClient) (string, error) {
	// prepare the base request
	req := identity.ListCompartmentsRequest{
		CompartmentId:          &tenancyOCID,
		AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
		LifecycleState:         identity.CompartmentLifecycleStateActive,
		CompartmentIdInSubtree: common.Bool(true),
	}

	// paginate through results; stop when OpcNextPage is nil
	pageToken := ""
	for {
		if pageToken != "" {
			req.Page = common.String(pageToken)
		}

		resp, err := idClient.ListCompartments(ctx, req)
		if err != nil {
			return "", fmt.Errorf("listing compartments: %w", err)
		}

		// scan each compartment summary for a name match
		for _, comp := range resp.Items {
			if comp.Name != nil && *comp.Name == compartmentName {
				return *comp.Id, nil
			}
		}

		// if there's no next page, we’re done searching
		if resp.OpcNextPage == nil {
			break
		}
		pageToken = *resp.OpcNextPage
	}

	return "", fmt.Errorf("compartment %q not found under tenancy %s", compartmentName, tenancyOCID)
}
