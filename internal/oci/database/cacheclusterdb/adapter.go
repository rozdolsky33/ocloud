package cacheclusterdb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/redis"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter implements the domain.CacheClusterRepository interface for OCI.
type Adapter struct {
	redisClient   redis.RedisClusterClient
	networkClient core.VirtualNetworkClient
	subnetCache   map[string]*core.Subnet
	vcnCache      map[string]*core.Vcn
	nsgCache      map[string]*core.NetworkSecurityGroup
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	redisClient, err := redis.NewRedisClusterClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}
	netClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual network client: %w", err)
	}
	return &Adapter{
		redisClient:   redisClient,
		networkClient: netClient,
		subnetCache:   make(map[string]*core.Subnet),
		vcnCache:      make(map[string]*core.Vcn),
		nsgCache:      make(map[string]*core.NetworkSecurityGroup),
	}, nil
}

// GetCacheCluster retrieves a single OCI Cache (Redis) Cluster by ID and maps it to the domain model.
func (a *Adapter) GetCacheCluster(ctx context.Context, clusterId string) (*domain.CacheCluster, error) {
	response, err := a.redisClient.GetRedisCluster(ctx, redis.GetRedisClusterRequest{
		RedisClusterId: &clusterId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cache cluster %s: %w", clusterId, err)
	}

	cluster, err := a.enrichAndMapCacheCluster(ctx, response.RedisCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// ListCacheClusters retrieves a list of OCI Cache (Redis) clusters from OCI.
func (a *Adapter) ListCacheClusters(ctx context.Context, compartmentID string) ([]domain.CacheCluster, error) {
	var allClusters []domain.CacheCluster
	var page *string

	for {
		resp, err := a.redisClient.ListRedisClusters(ctx, redis.ListRedisClustersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list cache clusters: %w", err)
		}

		for _, item := range resp.Items {
			attrs := mapping.NewCacheClusterAttributesFromOCIRedisClusterSummary(item)
			allClusters = append(allClusters, *mapping.NewDomainCacheClusterFromAttrs(attrs))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return allClusters, nil
}

// ListEnrichedCacheClusters retrieves a list of OCI Cache clusters from OCI and enriches them.
// It fetches full details for each cluster including network information.
func (a *Adapter) ListEnrichedCacheClusters(ctx context.Context, compartmentID string) ([]domain.CacheCluster, error) {
	var results []domain.CacheCluster
	var page *string

	// First, list all cluster IDs
	var clusterIDs []string
	for {
		resp, err := a.redisClient.ListRedisClusters(ctx, redis.ListRedisClustersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list cache clusters: %w", err)
		}

		for _, item := range resp.Items {
			if item.Id != nil {
				clusterIDs = append(clusterIDs, *item.Id)
			}
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	// Now fetch full details for each cluster
	for _, clusterID := range clusterIDs {
		cluster, err := a.GetCacheCluster(ctx, clusterID)
		if err != nil {
			// Log error but continue with other clusters
			continue
		}
		if cluster != nil {
			results = append(results, *cluster)
		}
	}

	return results, nil
}

// enrichNetworkNames resolves display names for subnet, VCN, and NSGs.
func (a *Adapter) enrichNetworkNames(ctx context.Context, c *domain.CacheCluster) error {
	if c.SubnetId != "" {
		if sub, err := a.getSubnet(ctx, c.SubnetId); err == nil && sub != nil {
			if sub.DisplayName != nil {
				c.SubnetName = *sub.DisplayName
			}
			if sub.VcnId != nil {
				c.VcnID = *sub.VcnId
				if vcn, err := a.getVcn(ctx, *sub.VcnId); err == nil && vcn != nil && vcn.DisplayName != nil {
					c.VcnName = *vcn.DisplayName
				}
			}
		}
	}

	// Enrich NSG names
	if len(c.NsgIds) > 0 {
		var names []string
		for _, id := range c.NsgIds {
			if nsg, err := a.getNsg(ctx, id); err == nil && nsg != nil && nsg.DisplayName != nil {
				names = append(names, *nsg.DisplayName)
			}
		}
		c.NsgNames = names
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

// enrichDomainCacheCluster applies additional lookups (e.g., network names) to the mapped domain model.
func (a *Adapter) enrichDomainCacheCluster(ctx context.Context, c *domain.CacheCluster) error {
	return a.enrichNetworkNames(ctx, c)
}

// enrichAndMapCacheCluster maps a full OCI RedisCluster and enriches it.
func (a *Adapter) enrichAndMapCacheCluster(ctx context.Context, cluster redis.RedisCluster) (*domain.CacheCluster, error) {
	attrs := mapping.NewCacheClusterAttributesFromOCIRedisCluster(cluster)
	c := mapping.NewDomainCacheClusterFromAttrs(attrs)
	if err := a.enrichDomainCacheCluster(ctx, c); err != nil {
		return c, fmt.Errorf("enriching cache cluster %s: %w", c.ID, err)
	}
	return c, nil
}
