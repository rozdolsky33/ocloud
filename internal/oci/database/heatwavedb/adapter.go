package heatwavedb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/mysql"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter implements the domain.HeatWaveDatabaseRepository interface for OCI.
type Adapter struct {
	mysqlClient   mysql.DbSystemClient
	networkClient core.VirtualNetworkClient
	subnetCache   map[string]*core.Subnet
	vcnCache      map[string]*core.Vcn
	nsgCache      map[string]*core.NetworkSecurityGroup
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	mysqlClient, err := mysql.NewDbSystemClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL client: %w", err)
	}
	netClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual network client: %w", err)
	}
	return &Adapter{
		mysqlClient:   mysqlClient,
		networkClient: netClient,
		subnetCache:   make(map[string]*core.Subnet),
		vcnCache:      make(map[string]*core.Vcn),
		nsgCache:      make(map[string]*core.NetworkSecurityGroup),
	}, nil
}

// GetHeatWaveDatabase retrieves a single HeatWave Database and maps it to the domain model.
func (a *Adapter) GetHeatWaveDatabase(ctx context.Context, ocid string) (*domain.HeatWaveDatabase, error) {
	response, err := a.mysqlClient.GetDbSystem(ctx, mysql.GetDbSystemRequest{
		DbSystemId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get HeatWave database: %w", err)
	}
	db, err := a.enrichAndMapHeatWaveDatabase(ctx, response.DbSystem)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// ListHeatWaveDatabases retrieves a list of HeatWave databases from OCI.
func (a *Adapter) ListHeatWaveDatabases(ctx context.Context, compartmentID string) ([]domain.HeatWaveDatabase, error) {
	var allDatabases []domain.HeatWaveDatabase
	var page *string
	for {
		resp, err := a.mysqlClient.ListDbSystems(ctx, mysql.ListDbSystemsRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list HeatWave databases: %w", err)
		}
		for _, item := range resp.Items {
			allDatabases = append(allDatabases, *mapping.NewDomainHeatWaveDatabaseFromAttrs(mapping.NewHeatWaveDatabaseAttributesFromOCIDbSystemSummary(item)))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return allDatabases, nil
}

// ListEnrichedHeatWaveDatabases retrieves a list of HeatWave databases from OCI and enriches them.
// It fetches full DbSystem details for each database to get complete information including
// SubnetId, NsgIds, DataStorageSizeInGBs, ConfigurationId, and Maintenance details.
func (a *Adapter) ListEnrichedHeatWaveDatabases(ctx context.Context, compartmentID string) ([]domain.HeatWaveDatabase, error) {
	var results []domain.HeatWaveDatabase
	var page *string

	// First, list all database IDs
	var dbIDs []string
	for {
		resp, err := a.mysqlClient.ListDbSystems(ctx, mysql.ListDbSystemsRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list HeatWave databases: %w", err)
		}
		for _, item := range resp.Items {
			if item.Id != nil {
				dbIDs = append(dbIDs, *item.Id)
			}
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	// Now fetch full details for each database
	for _, dbID := range dbIDs {
		db, err := a.GetHeatWaveDatabase(ctx, dbID)
		if err != nil {
			// Log error but continue with other databases
			continue
		}
		if db != nil {
			results = append(results, *db)
		}
	}

	return results, nil
}

// enrichNetworkNames resolves display names for subnet, VCN, and NSGs.
func (a *Adapter) enrichNetworkNames(ctx context.Context, d *domain.HeatWaveDatabase) error {
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

// enrichDomainHeatWaveDB applies additional lookups (e.g., network names) to the mapped domain model.
func (a *Adapter) enrichDomainHeatWaveDB(ctx context.Context, d *domain.HeatWaveDatabase) error {
	return a.enrichNetworkNames(ctx, d)
}

// enrichAndMapHeatWaveDatabase maps a full OCI DbSystem and enriches it.
func (a *Adapter) enrichAndMapHeatWaveDatabase(ctx context.Context, full mysql.DbSystem) (*domain.HeatWaveDatabase, error) {
	d := mapping.NewDomainHeatWaveDatabaseFromAttrs(mapping.NewHeatWaveDatabaseAttributesFromOCIDbSystem(full))
	if err := a.enrichDomainHeatWaveDB(ctx, d); err != nil {
		return d, fmt.Errorf("enriching HeatWave database %s: %w", d.ID, err)
	}
	return d, nil
}

// enrichAndMapHeatWaveDatabasesFromSummaries maps summaries and enriches them (best-effort names).
func (a *Adapter) enrichAndMapHeatWaveDatabasesFromSummaries(ctx context.Context, items []mysql.DbSystemSummary) ([]domain.HeatWaveDatabase, error) {
	res := make([]domain.HeatWaveDatabase, 0, len(items))
	for _, it := range items {
		d := mapping.NewDomainHeatWaveDatabaseFromAttrs(mapping.NewHeatWaveDatabaseAttributesFromOCIDbSystemSummary(it))
		if err := a.enrichDomainHeatWaveDB(ctx, d); err != nil {
			return nil, fmt.Errorf("enriching HeatWave database %s: %w", d.ID, err)
		}
		res = append(res, *d)
	}
	return res, nil
}
