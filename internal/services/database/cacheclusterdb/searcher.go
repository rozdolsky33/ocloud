package cacheclusterdb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableCacheCluster adapts CacheCluster to the search.Indexable interface.
type SearchableCacheCluster struct {
	database.CacheCluster
}

// ToIndexable converts a CacheCluster to a map of searchable fields.
func (s SearchableCacheCluster) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	// Format node count
	var nodeCount string
	if s.NodeCount > 0 {
		nodeCount = strconv.Itoa(s.NodeCount)
	}

	// Format shard count
	var shardCount string
	if s.ShardCount > 0 {
		shardCount = strconv.Itoa(s.ShardCount)
	}

	// Format memory
	var memory string
	if s.NodeMemoryInGBs > 0 {
		memory = fmt.Sprintf("%.0f", s.NodeMemoryInGBs)
	}

	// join slices safely
	join := func(items []string) string {
		return strings.ToLower(strings.Join(items, ","))
	}

	return map[string]any{
		"ID":                         strings.ToLower(s.ID),
		"DisplayName":                strings.ToLower(s.DisplayName),
		"State":                      strings.ToLower(s.LifecycleState),
		"SoftwareVersion":            strings.ToLower(s.SoftwareVersion),
		"ClusterMode":                strings.ToLower(s.ClusterMode),
		"NodeCount":                  strings.ToLower(nodeCount),
		"ShardCount":                 strings.ToLower(shardCount),
		"NodeMemoryInGBs":            strings.ToLower(memory),
		"PrimaryFqdn":                strings.ToLower(s.PrimaryFqdn),
		"PrimaryEndpointIpAddress":   strings.ToLower(s.PrimaryEndpointIpAddress),
		"ReplicasFqdn":               strings.ToLower(s.ReplicasFqdn),
		"ReplicasEndpointIpAddress":  strings.ToLower(s.ReplicasEndpointIpAddress),
		"DiscoveryFqdn":              strings.ToLower(s.DiscoveryFqdn),
		"DiscoveryEndpointIpAddress": strings.ToLower(s.DiscoveryEndpointIpAddress),
		"VcnID":                      strings.ToLower(s.VcnID),
		"VcnName":                    strings.ToLower(s.VcnName),
		"SubnetId":                   strings.ToLower(s.SubnetId),
		"SubnetName":                 strings.ToLower(s.SubnetName),
		"NsgNames":                   join(s.NsgNames),
		"NsgIds":                     join(s.NsgIds),
		"TagsKV":                     strings.ToLower(tagsKV),
		"TagsVal":                    strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the list of fields to be indexed for Cache Clusters.
func GetSearchableFields() []string {
	return []string{
		"ID", "DisplayName", "State", "SoftwareVersion", "ClusterMode",
		"NodeCount", "ShardCount", "NodeMemoryInGBs",
		"PrimaryFqdn", "PrimaryEndpointIpAddress",
		"ReplicasFqdn", "ReplicasEndpointIpAddress",
		"DiscoveryFqdn", "DiscoveryEndpointIpAddress",
		"VcnID", "VcnName", "SubnetId", "SubnetName",
		"NsgNames", "NsgIds",
		"TagsKV", "TagsVal",
	}
}

// GetBoostedFields returns the list of fields to be boosted in the search.
func GetBoostedFields() []string {
	return []string{"DisplayName", "ID", "VcnName", "SubnetName", "PrimaryFqdn"}
}

// ToSearchableCacheClusters converts a slice of CacheCluster to a slice of search.Indexable.
func ToSearchableCacheClusters(clusters []database.CacheCluster) []search.Indexable {
	searchable := make([]search.Indexable, len(clusters))
	for i, cluster := range clusters {
		searchable[i] = SearchableCacheCluster{cluster}
	}
	return searchable
}
