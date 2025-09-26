package autonomousdb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/database"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter implements the domain.AutonomousDatabaseRepository interface for OCI.
type Adapter struct {
	dbClient      database.DatabaseClient
	networkClient core.VirtualNetworkClient
	subnetCache   map[string]*core.Subnet
	vcnCache      map[string]*core.Vcn
	nsgCache      map[string]*core.NetworkSecurityGroup
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	dbClient, err := oci.NewDatabaseClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	netClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual network client: %w", err)
	}
	return &Adapter{
		dbClient:      dbClient,
		networkClient: netClient,
		subnetCache:   make(map[string]*core.Subnet),
		vcnCache:      make(map[string]*core.Vcn),
		nsgCache:      make(map[string]*core.NetworkSecurityGroup),
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
	db, err := a.enrichAndMapAutonomousDatabase(ctx, response.AutonomousDatabase)
	if err != nil {
		return nil, err
	}
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

// ListEnrichedAutonomousDatabase retrieves a list of autonomous databases from OCI and enriches them.
func (a *Adapter) ListEnrichedAutonomousDatabase(ctx context.Context, compartmentID string) ([]domain.AutonomousDatabase, error) {
	var results []domain.AutonomousDatabase
	var page *string
	for {
		resp, err := a.dbClient.ListAutonomousDatabases(ctx, database.ListAutonomousDatabasesRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list autonomous databases: %w", err)
		}
		// Map summaries then enrich (network names). Keep it lightweight like instances' batch enrichment.
		batch, err := a.enrichAndMapAutonomousDatabasesFromSummaries(ctx, resp.Items)
		if err != nil {
			return nil, err
		}
		results = append(results, batch...)
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return results, nil
}

// enrichNetworkNames resolves display names for subnet, VCN, and NSGs.
func (a *Adapter) enrichNetworkNames(ctx context.Context, d *domain.AutonomousDatabase) error {
	if d.SubnetId != "" {
		if sub, err := a.getSubnet(ctx, d.SubnetId); err == nil && sub != nil {
			if sub.DisplayName != nil {
				d.SubnetName = *sub.DisplayName
			}
			if sub.VcnId != nil {
				d.VcnID = *sub.VcnId
				if vcn, err := a.getVcn(ctx, *sub.VcnId); err == nil && vcn != nil && vcn.DisplayName != nil {
					d.VcnName = *vcn.DisplayName
				}
			}
		}
	}
	if len(d.NsgIds) > 0 {
		var names []string
		for _, id := range d.NsgIds {
			if nsg, err := a.getNsg(ctx, id); err == nil && nsg != nil && nsg.DisplayName != nil {
				names = append(names, *nsg.DisplayName)
			}
		}
		d.NsgNames = names
	}
	return nil
}

// getSubnet retrieves a subnet by its ID, utilizing a local cache for improved performance.
func (a *Adapter) getSubnet(ctx context.Context, id string) (*core.Subnet, error) {
	if s, ok := a.subnetCache[id]; ok {
		return s, nil
	}
	resp, err := a.networkClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &id})
	if err != nil {
		return nil, err
	}
	a.subnetCache[id] = &resp.Subnet
	return &resp.Subnet, nil
}

// getVcn retrieves a VCN by its ID, utilizing a local cache for improved performance.
func (a *Adapter) getVcn(ctx context.Context, id string) (*core.Vcn, error) {
	if v, ok := a.vcnCache[id]; ok {
		return v, nil
	}
	resp, err := a.networkClient.GetVcn(ctx, core.GetVcnRequest{VcnId: &id})
	if err != nil {
		return nil, err
	}
	a.vcnCache[id] = &resp.Vcn
	return &resp.Vcn, nil
}

// getNsg retrieves a NSG by its ID, utilizing a local cache for improved performance.
func (a *Adapter) getNsg(ctx context.Context, id string) (*core.NetworkSecurityGroup, error) {
	if n, ok := a.nsgCache[id]; ok {
		return n, nil
	}
	resp, err := a.networkClient.GetNetworkSecurityGroup(ctx, core.GetNetworkSecurityGroupRequest{NetworkSecurityGroupId: &id})
	if err != nil {
		return nil, err
	}
	a.nsgCache[id] = &resp.NetworkSecurityGroup
	return &resp.NetworkSecurityGroup, nil
}

// enrichDomainAutonomousDB applies additional lookups (e.g., network names) to the mapped domain model.
func (a *Adapter) enrichDomainAutonomousDB(ctx context.Context, d *domain.AutonomousDatabase) error {
	return a.enrichNetworkNames(ctx, d)
}

// enrichAndMapAutonomousDatabase maps a full OCI AutonomousDatabase and enriches it.
func (a *Adapter) enrichAndMapAutonomousDatabase(ctx context.Context, full database.AutonomousDatabase) (domain.AutonomousDatabase, error) {
	d := a.toDomainAutonomousDB(full)
	if err := a.enrichDomainAutonomousDB(ctx, &d); err != nil {
		return d, fmt.Errorf("enriching autonomous database %s: %w", d.ID, err)
	}
	return d, nil
}

// enrichAndMapAutonomousDatabasesFromSummaries maps summaries and enriches them (best-effort names).
func (a *Adapter) enrichAndMapAutonomousDatabasesFromSummaries(ctx context.Context, items []database.AutonomousDatabaseSummary) ([]domain.AutonomousDatabase, error) {
	res := make([]domain.AutonomousDatabase, 0, len(items))
	for _, it := range items {
		d := a.toDomainAutonomousDB(it)
		if err := a.enrichDomainAutonomousDB(ctx, &d); err != nil {
			return nil, fmt.Errorf("enriching autonomous database %s: %w", d.ID, err)
		}
		res = append(res, d)
	}
	return res, nil
}

// toDomainAutonomousDB maps either a full database.AutonomousDatabase (from Get) or a database.AutonomousDatabaseSummary (from List) into the single domain.AutonomousDatabase type.
func (a *Adapter) toDomainAutonomousDB(ociObj interface{}) domain.AutonomousDatabase {
	var (
		// identity
		name            *string
		id              *string
		compartmentId   *string
		lifecycleState  string
		dbVersion       *string
		dbWorkloadStr   string
		licenseModelStr string

		// network
		whitelistedIps       []string
		privateEndpoint      *string
		privateEndpointIp    *string
		privateEndpointLabel *string
		subnetId             *string
		nsgIds               []string
		isMtlsRequired       *bool
		isPubliclyAccessible *bool

		// capacity
		computeModelStr             string
		ecpuCount                   *float32
		ocpuCount                   *float32
		cpuCoreCount                *int
		storageTBs                  *int
		storageGBs                  *int
		isAutoScale                 *bool
		isStorageAutoScalingEnabled *bool

		// operations & integrations
		operationsInsightsStatus string
		databaseManagementStatus string
		dataSafeStatus           string
		isFreeTier               *bool

		// Data Guard / DR
		isDataGuardEnabled  *bool
		role                *string
		peerAutonomousDbIds []string

		// maintenance
		patchModel           *string
		nextMaintenanceRunId *string
		maintenanceSchedule  *string

		// connections
		connStrings *database.AutonomousDatabaseConnectionStrings
		connUrls    *database.AutonomousDatabaseConnectionUrls

		freeformTags map[string]string
		definedTags  map[string]map[string]interface{}

		timeCreated *common.SDKTime
	)

	switch src := ociObj.(type) {
	case database.AutonomousDatabase:
		name = src.DbName
		id = src.Id
		compartmentId = src.CompartmentId
		lifecycleState = string(src.LifecycleState)
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

		computeModelStr = string(src.ComputeModel)
		ecpuCount = src.ComputeCount
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

		computeModelStr = string(src.ComputeModel)
		ecpuCount = src.ComputeCount
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
	switch src := ociObj.(type) {
	case database.AutonomousDatabase:
		storageGBs = src.DataStorageSizeInGBs
	case database.AutonomousDatabaseSummary:
		storageGBs = src.DataStorageSizeInGBs
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
	if dbVersion != nil {
		d.DbVersion = *dbVersion
	}
	d.DbWorkload = dbWorkloadStr
	d.LicenseModel = licenseModelStr

	// network
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
	d.ComputeModel = computeModelStr
	d.EcpuCount = ecpuCount
	d.OcpuCount = ocpuCount
	d.CpuCoreCount = cpuCoreCount
	d.DataStorageSizeInTBs = storageTBs
	d.DataStorageSizeInGBs = storageGBs
	d.IsAutoScalingEnabled = isAutoScale
	d.IsStorageAutoScalingEnabled = isStorageAutoScalingEnabled
	d.OperationsInsightsStatus = operationsInsightsStatus
	d.DatabaseManagementStatus = databaseManagementStatus
	d.DataSafeStatus = dataSafeStatus
	d.IsFreeTier = isFreeTier

	d.IsPubliclyAccessible = isPubliclyAccessible

	d.IsDataGuardEnabled = isDataGuardEnabled
	if role != nil {
		d.Role = *role
	}
	d.PeerAutonomousDbIds = peerAutonomousDbIds

	if patchModel != nil {
		d.PatchModel = *patchModel
	}
	if nextMaintenanceRunId != nil {
		d.NextMaintenanceRunId = *nextMaintenanceRunId
	}
	if maintenanceSchedule != nil {
		d.MaintenanceScheduleType = *maintenanceSchedule
	}

	if connStrings != nil {
		if connStrings.AllConnectionStrings != nil {
			d.ConnectionStrings = connStrings.AllConnectionStrings
		}
		d.Profiles = connStrings.Profiles
	}
	d.ConnectionUrls = connUrls
	d.FreeformTags = freeformTags
	d.DefinedTags = definedTags
	if timeCreated != nil {
		t := timeCreated.Time
		d.TimeCreated = &t
	}

	return d
}
