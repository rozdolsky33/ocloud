package cacheclusterdb

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewCacheClusterListModel builds a TUI list for OCI Cache (Redis/Valkey) Clusters.
func NewCacheClusterListModel(clusters []domain.CacheCluster) tui.Model {
	return tui.NewModel("OCI Cache Clusters", clusters, func(cluster domain.CacheCluster) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          cluster.ID,
			Title:       cluster.DisplayName,
			Description: describeCacheCluster(cluster),
		}
	})
}

func describeCacheCluster(cluster domain.CacheCluster) string {
	// Node count and memory
	resourceInfo := ""
	if cluster.NodeCount > 0 {
		resourceInfo = fmt.Sprintf("%d nodes", cluster.NodeCount)
		if cluster.NodeMemoryInGBs > 0 {
			resourceInfo = fmt.Sprintf("%d nodes × %.0fGB", cluster.NodeCount, cluster.NodeMemoryInGBs)
		}
	}

	// Cluster mode (sharded/non-sharded)
	mode := ""
	if cluster.ClusterMode != "" {
		mode = cluster.ClusterMode
		if cluster.ClusterMode == "SHARDED" && cluster.ShardCount > 0 {
			mode = fmt.Sprintf("%s (%d shards)", cluster.ClusterMode, cluster.ShardCount)
		}
	}

	// Software version (Redis/Valkey)
	version := ""
	if cluster.SoftwareVersion != "" {
		version = cluster.SoftwareVersion
	}

	// Network - subnet name or VCN name
	network := ""
	if cluster.SubnetName != "" {
		network = cluster.SubnetName
	} else if cluster.VcnName != "" {
		network = cluster.VcnName
	}

	// Date created
	date := ""
	if cluster.TimeCreated != nil && !cluster.TimeCreated.IsZero() {
		date = cluster.TimeCreated.Format("2006-01-02")
	}

	// Build description parts
	parts := []string{}
	if cluster.LifecycleState != "" {
		parts = append(parts, cluster.LifecycleState)
	}
	if resourceInfo != "" {
		parts = append(parts, resourceInfo)
	}
	if mode != "" {
		parts = append(parts, mode)
	}
	if version != "" {
		parts = append(parts, version)
	}
	if network != "" {
		parts = append(parts, network)
	}
	if date != "" {
		parts = append(parts, date)
	}

	return strings.Join(parts, " • ")
}
