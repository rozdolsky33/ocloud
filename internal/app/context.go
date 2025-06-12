package app

import (
	"context"
	"fmt"
	"os"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/logger"
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

// NewAppContext initializes AppContext, resolves tenancy & compartment IDs, and builds OCI clients.
func NewAppContext(ctx context.Context, cmd *cobra.Command) (*AppContext, error) {
	log := logger.CmdLogger
	log.Info("Initializing application context")

	// Load OCI config & create provider
	prov := config.LoadOCIConfig()

	// Create an identity client (needed for compartment lookup)
	idClient, err := identity.NewIdentityClientWithConfigurationProvider(prov)
	if err != nil {
		return nil, fmt.Errorf("creating identity client: %w", err)
	}

	// Optional region override
	if region, ok := os.LookupEnv(config.EnvOCIRegion); ok {
		idClient.SetRegion(region)
		log.V(1).Info("overriding region from env", "region", region)
	}

	// Build base AppContext
	appCtx := &AppContext{
		Ctx:             ctx,
		Provider:        prov,
		IdentityClient:  idClient,
		CompartmentName: viper.GetString(config.FlagNameCompartment),
	}

	// Resolve Tenancy ID
	tenancyID, err := ResolveTenancyID(cmd)
	if err != nil {
		return nil, err
	}
	appCtx.TenancyID = tenancyID

	// Resolve Tenancy Name
	tenancyName := ResolveTenancyName(cmd, tenancyID)
	if tenancyName != "" {
		appCtx.TenancyName = tenancyName
	}

	// Resolve Compartment ID
	compID, err := ResolveCompartmentID(ctx, appCtx.TenancyID, appCtx.CompartmentName, appCtx.IdentityClient)
	if err != nil {
		return nil, fmt.Errorf("could not resolve compartment ID: %w", err)
	}
	appCtx.CompartmentID = compID

	return appCtx, nil
}

// ResolveTenancyID resolves the tenancy OCID from various sources in order of precedence:
// 1. Command line flag
// 2. Environment variable
// 3. Tenancy name lookup (if tenancy name is provided)
// 4. OCI config file
// Returns the tenancy ID or an error if it cannot be resolved.
func ResolveTenancyID(cmd *cobra.Command) (string, error) {
	log := logger.CmdLogger

	// Check if tenancy ID is provided as a flag
	if cmd.Flags().Changed(config.FlagNameTenancyID) {
		tenancyID := viper.GetString(config.FlagNameTenancyID)
		log.V(1).Info("using tenancy OCID from flag", "tenancyID", tenancyID)
		return tenancyID, nil
	}

	// Check if tenancy ID is provided as an environment variable
	if envTenancy := os.Getenv(config.EnvOCITenancy); envTenancy != "" {
		log.V(1).Info("using tenancy OCID from env", "tenancyID", envTenancy)
		viper.Set(config.FlagNameTenancyID, envTenancy)
		return envTenancy, nil
	}

	// Check if the tenancy name is provided as an environment variable
	if envTenancyName := os.Getenv(config.EnvOCITenancyName); envTenancyName != "" {
		lookupID, err := config.LookupTenancyID(envTenancyName)
		if err != nil {
			// Log the error but continue with the next method of resolving the tenancy ID
			log.Info("could not look up tenancy ID for tenancy name, continuing with other methods", "tenancyName", envTenancyName, "error", err)
			// Add a more detailed message about how to set up the mapping file
			log.Info("To set up tenancy mapping, create a YAML file at ~/.oci/tenancy-map.yaml or set the OCI_TENANCY_MAP_PATH environment variable. The file should contain entries mapping tenancy names to OCIDs. Example:\n- environment: prod\n  tenancy: mytenancy\n  tenancy_id: ocid1.tenancy.oc1..aaaaaaaabcdefghijklmnopqrstuvwxyz\n  realm: oc1\n  compartments: mycompartment\n  regions: us-ashburn-1")
		} else {
			log.V(1).Info("using tenancy OCID for name", "tenancyName", envTenancyName, "tenancyID", lookupID)
			viper.Set(config.FlagNameTenancyID, lookupID)
			return lookupID, nil
		}
	}

	// Load from an OCI config file as a last resort
	tenancyID, err := config.GetTenancyOCID()
	if err != nil {
		return "", fmt.Errorf("could not load tenancy OCID: %w", err)
	}
	log.V(1).Info("using tenancy OCID from config file", "tenancyID", tenancyID)
	viper.Set(config.FlagNameTenancyID, tenancyID)

	return tenancyID, nil
}

// ResolveTenancyName resolves the tenancy name from various sources in order of precedence:
// 1. Command line flag
// 2. Environment variable
// 3. Tenancy mapping file lookup (using tenancy ID)
// Returns the tenancy name or an empty string if it cannot be resolved.
func ResolveTenancyName(cmd *cobra.Command, tenancyID string) string {
	log := logger.CmdLogger

	// Check if the tenancy name is provided as a flag
	if cmd.Flags().Changed(config.FlagNameTenancyName) {
		tenancyName := viper.GetString(config.FlagNameTenancyName)
		log.V(1).Info("using tenancy name from flag", "tenancyName", tenancyName)
		return tenancyName
	}

	// Check if the tenancy name is provided as an environment variable
	if envTenancyName := os.Getenv(config.EnvOCITenancyName); envTenancyName != "" {
		log.V(1).Info("using tenancy name from env", "tenancyName", envTenancyName)
		viper.Set(config.FlagNameTenancyName, envTenancyName)
		return envTenancyName
	}

	// Try to find a tenancy name from a mapping file if available
	tenancies, err := config.LoadTenancyMap()
	if err == nil {
		for _, env := range tenancies {
			if env.TenancyID == tenancyID {
				log.V(1).Info("found tenancy name from mapping file", "tenancyName", env.Tenancy)
				viper.Set(config.FlagNameTenancyName, env.Tenancy)
				return env.Tenancy
			}
		}
	}

	return ""
}

// ResolveCompartmentID returns the OCID of the compartment whose name matches
// `compartmentName` under the given tenancy. It searches all active compartments
// in the tenancy subtree.
func ResolveCompartmentID(ctx context.Context, tenancyOCID, compartmentName string, idClient identity.IdentityClient) (string, error) {
	// If the compartment name is not set, use tenancy ID as fallback
	if compartmentName == "" {
		logger.CmdLogger.V(1).Info("compartment name not set, using tenancy ID as fallback", "tenancyID", tenancyOCID)
		return tenancyOCID, nil
	}

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

		// if there's no next page, we're done searching
		if resp.OpcNextPage == nil {
			break
		}
		pageToken = *resp.OpcNextPage
	}

	return "", fmt.Errorf("compartment %q not found under tenancy %s", compartmentName, tenancyOCID)
}
