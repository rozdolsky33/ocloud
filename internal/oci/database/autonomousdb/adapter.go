package autonomousdb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter implements the domain.AutonomousDatabaseRepository interface for OCI.
type Adapter struct {
	dbClient      database.DatabaseClient
	networkClient core.VirtualNetworkClient
	// simple caches to avoid repeated lookups during a run
	subnetCache map[string]*core.Subnet
	vcnCache    map[string]*core.Vcn
	nsgCache    map[string]*core.NetworkSecurityGroup
}

// NewAdapter creates a new Adapter instance. The compartmentID parameter is accepted for
// backward compatibility with service wiring but is not required by the adapter itself.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	dbClient, err := oci.NewDatabaseClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	// create a virtual network client for name enrichment
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
	db := a.toDomainAutonomousDB(response.AutonomousDatabase)
	// Best-effort network name enrichment (subnet, VCN, NSGs)
	if err := a.enrichNetworkNames(ctx, &db); err != nil {
		// non-fatal; keep basic info if lookups fail
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

// enrichNetworkNames resolves display names for subnet, VCN, and NSGs where possible (best effort, cached).
func (a *Adapter) enrichNetworkNames(ctx context.Context, d *domain.AutonomousDatabase) error {
	// Subnet → name and VCN
	if d.SubnetId != "" {
		// subnet
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
	// NSGs → names
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

// cached lookups
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

	// Post-switch normalization for fields not captured in shared vars
	switch src := ociObj.(type) {
	case database.AutonomousDatabase:
		// Capacity (additional fields)
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
	d.ComputeModel = computeModelStr
	d.EcpuCount = ecpuCount
	d.OcpuCount = ocpuCount
	d.CpuCoreCount = cpuCoreCount
	d.DataStorageSizeInTBs = storageTBs
	d.DataStorageSizeInGBs = storageGBs
	d.IsAutoScalingEnabled = isAutoScale
	// additional capacity flags
	d.IsStorageAutoScalingEnabled = isStorageAutoScalingEnabled

	// operations & integrations
	d.OperationsInsightsStatus = operationsInsightsStatus
	d.DatabaseManagementStatus = databaseManagementStatus
	d.DataSafeStatus = dataSafeStatus
	d.IsFreeTier = isFreeTier

	// networking extras
	d.IsPubliclyAccessible = isPubliclyAccessible

	// Data Guard / DR
	d.IsDataGuardEnabled = isDataGuardEnabled
	if role != nil {
		d.Role = *role
	}
	d.PeerAutonomousDbIds = peerAutonomousDbIds

	// maintenance
	if patchModel != nil {
		d.PatchModel = *patchModel
	}
	if nextMaintenanceRunId != nil {
		d.NextMaintenanceRunId = *nextMaintenanceRunId
	}
	if maintenanceSchedule != nil {
		d.MaintenanceScheduleType = *maintenanceSchedule
	}

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
