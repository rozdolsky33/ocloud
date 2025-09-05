package autonomousdb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter implements the domain.AutonomousDatabaseRepository interface for OCI.
type Adapter struct {
	dbClient database.DatabaseClient
}

// NewAdapter creates a new Adapter instance. The compartmentID parameter is accepted for
// backward compatibility with service wiring but is not required by the adapter itself.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	dbClient, err := oci.NewDatabaseClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	return &Adapter{
		dbClient: dbClient,
	}, nil
}

// GetAutonomousDatabase retrieves a single Autonomous Database and maps it to the domain model.
func (a *Adapter) GetAutonomousDatabase(ctx context.Context, ocid string) (*domain.AutonomousDatabase, error) {
	response, err := a.dbClient.GetAutonomousDatabase(ctx, database.GetAutonomousDatabaseRequest{
		AutonomousDatabaseId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get autonomous database: %w", err)
	}
	db := a.toDomainAutonomousDB(response.AutonomousDatabase)
	return &db, nil
}

// ListAutonomousDatabases retrieves a list of autonomous databases from OCI.
func (a *Adapter) ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]domain.AutonomousDatabase, error) {
	var allDatabases []domain.AutonomousDatabase
	var page *string
	for {
		resp, err := a.dbClient.ListAutonomousDatabases(ctx, database.ListAutonomousDatabasesRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list autonomous databases: %w", err)
		}
		for _, item := range resp.Items {
			allDatabases = append(allDatabases, a.toDomainAutonomousDB(item))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return allDatabases, nil
}

// toDomainAutonomousDB maps either a full database.AutonomousDatabase (from Get) or a database.AutonomousDatabaseSummary (from List) into the single domain.AutonomousDatabase type.
func (a *Adapter) toDomainAutonomousDB(ociObj interface{}) domain.AutonomousDatabase {
	var (
		// identity
		name             *string
		id               *string
		compartmentId    *string
		lifecycleState   string
		lifecycleDetails *string
		dbVersion        *string
		dbWorkloadStr    string
		licenseModelStr  string

		// networking
		whitelistedIps       []string
		privateEndpoint      *string
		privateEndpointIp    *string
		privateEndpointLabel *string
		subnetId             *string
		nsgIds               []string
		isMtlsRequired       *bool

		// capacity
		ocpuCount    *float32
		cpuCoreCount *int
		storageTBs   *int
		isAutoScale  *bool

		// connections
		connStrings *database.AutonomousDatabaseConnectionStrings
		connUrls    *database.AutonomousDatabaseConnectionUrls

		// tags
		freeformTags map[string]string
		definedTags  map[string]map[string]interface{}

		// timestamps
		timeCreated *common.SDKTime
	)

	switch src := ociObj.(type) {
	case database.AutonomousDatabase:
		name = src.DbName
		id = src.Id
		compartmentId = src.CompartmentId
		lifecycleState = string(src.LifecycleState)
		lifecycleDetails = src.LifecycleDetails
		dbVersion = src.DbVersion
		dbWorkloadStr = string(src.DbWorkload)
		licenseModelStr = string(src.LicenseModel)

		whitelistedIps = src.WhitelistedIps
		privateEndpoint = src.PrivateEndpoint
		privateEndpointIp = src.PrivateEndpointIp
		privateEndpointLabel = src.PrivateEndpointLabel
		subnetId = src.SubnetId
		nsgIds = src.NsgIds
		isMtlsRequired = src.IsMtlsConnectionRequired

		ocpuCount = src.OcpuCount
		cpuCoreCount = src.CpuCoreCount
		storageTBs = src.DataStorageSizeInTBs
		isAutoScale = src.IsAutoScalingEnabled

		connStrings = src.ConnectionStrings
		connUrls = src.ConnectionUrls

		freeformTags = src.FreeformTags
		definedTags = src.DefinedTags
		timeCreated = src.TimeCreated

	case database.AutonomousDatabaseSummary:
		name = src.DbName
		id = src.Id
		compartmentId = src.CompartmentId
		lifecycleState = string(src.LifecycleState)
		lifecycleDetails = src.LifecycleDetails
		dbVersion = src.DbVersion
		dbWorkloadStr = string(src.DbWorkload)
		licenseModelStr = string(src.LicenseModel)

		whitelistedIps = src.WhitelistedIps
		privateEndpoint = src.PrivateEndpoint
		privateEndpointIp = src.PrivateEndpointIp
		privateEndpointLabel = src.PrivateEndpointLabel
		subnetId = src.SubnetId
		nsgIds = src.NsgIds
		isMtlsRequired = src.IsMtlsConnectionRequired

		ocpuCount = src.OcpuCount
		cpuCoreCount = src.CpuCoreCount
		storageTBs = src.DataStorageSizeInTBs
		isAutoScale = src.IsAutoScalingEnabled

		connStrings = src.ConnectionStrings
		connUrls = src.ConnectionUrls

		freeformTags = src.FreeformTags
		definedTags = src.DefinedTags
		timeCreated = src.TimeCreated
	default:
		return domain.AutonomousDatabase{}
	}

	d := domain.AutonomousDatabase{}
	if name != nil {
		d.Name = *name
	}
	if id != nil {
		d.ID = *id
	}
	if compartmentId != nil {
		d.CompartmentOCID = *compartmentId
	}
	d.LifecycleState = lifecycleState
	if lifecycleDetails != nil {
		d.LifecycleDetails = *lifecycleDetails
	}
	if dbVersion != nil {
		d.DbVersion = *dbVersion
	}
	d.DbWorkload = dbWorkloadStr
	d.LicenseModel = licenseModelStr

	// networking
	d.WhitelistedIps = whitelistedIps
	if privateEndpoint != nil {
		d.PrivateEndpoint = *privateEndpoint
	}
	if privateEndpointIp != nil {
		d.PrivateEndpointIp = *privateEndpointIp
	}
	if privateEndpointLabel != nil {
		d.PrivateEndpointLabel = *privateEndpointLabel
	}
	if subnetId != nil {
		d.SubnetId = *subnetId
	}
	d.NsgIds = nsgIds
	d.IsMtlsRequired = isMtlsRequired

	// capacity
	d.OcpuCount = ocpuCount
	d.CpuCoreCount = cpuCoreCount
	d.DataStorageSizeInTBs = storageTBs
	d.IsAutoScalingEnabled = isAutoScale

	// connections
	if connStrings != nil {
		if connStrings.AllConnectionStrings != nil {
			d.ConnectionStrings = connStrings.AllConnectionStrings
		}
		d.Profiles = connStrings.Profiles
	}
	d.ConnectionUrls = connUrls

	// tags
	d.FreeformTags = freeformTags
	d.DefinedTags = definedTags

	// timestamps
	if timeCreated != nil {
		t := timeCreated.Time
		d.TimeCreated = &t
	}

	return d
}
