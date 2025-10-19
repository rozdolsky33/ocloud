package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/redis"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
)

// CacheClusterAttributes holds intermediate attributes for mapping from OCI SDK to a domain model.
type CacheClusterAttributes struct {
	ID                         *string
	DisplayName                *string
	CompartmentOCID            *string
	LifecycleState             string
	LifecycleDetails           *string
	TimeCreated                *common.SDKTime
	TimeUpdated                *common.SDKTime
	NodeCount                  *int
	NodeMemoryInGBs            *float32
	SoftwareVersion            string
	ClusterMode                string
	ShardCount                 *int
	ConfigSetId                *string
	SubnetId                   *string
	NsgIds                     []string
	PrimaryFqdn                *string
	PrimaryEndpointIpAddress   *string
	ReplicasFqdn               *string
	ReplicasEndpointIpAddress  *string
	DiscoveryFqdn              *string
	DiscoveryEndpointIpAddress *string
	Nodes                      []redis.Node
	FreeformTags               map[string]string
	DefinedTags                map[string]map[string]interface{}
	SystemTags                 map[string]map[string]interface{}
}

// NewCacheClusterAttributesFromOCIRedisCluster converts a full OCI RedisCluster to attributes.
func NewCacheClusterAttributesFromOCIRedisCluster(cluster redis.RedisCluster) *CacheClusterAttributes {
	var nodes []redis.Node
	if cluster.NodeCollection != nil && cluster.NodeCollection.Items != nil {
		nodes = cluster.NodeCollection.Items
	}

	return &CacheClusterAttributes{
		ID:                         cluster.Id,
		DisplayName:                cluster.DisplayName,
		CompartmentOCID:            cluster.CompartmentId,
		LifecycleState:             string(cluster.LifecycleState),
		LifecycleDetails:           cluster.LifecycleDetails,
		TimeCreated:                cluster.TimeCreated,
		TimeUpdated:                cluster.TimeUpdated,
		NodeCount:                  cluster.NodeCount,
		NodeMemoryInGBs:            cluster.NodeMemoryInGBs,
		SoftwareVersion:            string(cluster.SoftwareVersion),
		ClusterMode:                string(cluster.ClusterMode),
		ShardCount:                 cluster.ShardCount,
		ConfigSetId:                cluster.OciCacheConfigSetId,
		SubnetId:                   cluster.SubnetId,
		NsgIds:                     cluster.NsgIds,
		PrimaryFqdn:                cluster.PrimaryFqdn,
		PrimaryEndpointIpAddress:   cluster.PrimaryEndpointIpAddress,
		ReplicasFqdn:               cluster.ReplicasFqdn,
		ReplicasEndpointIpAddress:  cluster.ReplicasEndpointIpAddress,
		DiscoveryFqdn:              cluster.DiscoveryFqdn,
		DiscoveryEndpointIpAddress: cluster.DiscoveryEndpointIpAddress,
		Nodes:                      nodes,
		FreeformTags:               cluster.FreeformTags,
		DefinedTags:                cluster.DefinedTags,
		SystemTags:                 cluster.SystemTags,
	}
}

// NewCacheClusterAttributesFromOCIRedisClusterSummary converts an OCI RedisClusterSummary to attributes.
func NewCacheClusterAttributesFromOCIRedisClusterSummary(summary redis.RedisClusterSummary) *CacheClusterAttributes {
	return &CacheClusterAttributes{
		ID:                         summary.Id,
		DisplayName:                summary.DisplayName,
		CompartmentOCID:            summary.CompartmentId,
		LifecycleState:             string(summary.LifecycleState),
		LifecycleDetails:           summary.LifecycleDetails,
		TimeCreated:                summary.TimeCreated,
		TimeUpdated:                summary.TimeUpdated,
		NodeCount:                  summary.NodeCount,
		NodeMemoryInGBs:            summary.NodeMemoryInGBs,
		SoftwareVersion:            string(summary.SoftwareVersion),
		ClusterMode:                string(summary.ClusterMode),
		ShardCount:                 summary.ShardCount,
		ConfigSetId:                summary.OciCacheConfigSetId,
		SubnetId:                   summary.SubnetId,
		NsgIds:                     summary.NsgIds,
		PrimaryFqdn:                summary.PrimaryFqdn,
		PrimaryEndpointIpAddress:   summary.PrimaryEndpointIpAddress,
		ReplicasFqdn:               summary.ReplicasFqdn,
		ReplicasEndpointIpAddress:  summary.ReplicasEndpointIpAddress,
		DiscoveryFqdn:              summary.DiscoveryFqdn,
		DiscoveryEndpointIpAddress: summary.DiscoveryEndpointIpAddress,
		Nodes:                      nil, // Not available in summary
		FreeformTags:               summary.FreeformTags,
		DefinedTags:                summary.DefinedTags,
		SystemTags:                 summary.SystemTags,
	}
}

// NewDomainCacheClusterFromAttrs converts CacheClusterAttributes to domain.CacheCluster.
func NewDomainCacheClusterFromAttrs(attrs *CacheClusterAttributes) *domain.CacheCluster {
	// Helper to safely dereference string pointers
	val := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	// Helper to safely dereference int pointers
	intVal := func(p *int) int {
		if p == nil {
			return 0
		}
		return *p
	}

	// Helper to safely dereference float32 pointers
	float32Val := func(p *float32) float32 {
		if p == nil {
			return 0
		}
		return *p
	}

	var timeCreated, timeUpdated *time.Time
	if attrs.TimeCreated != nil {
		t := attrs.TimeCreated.Time
		timeCreated = &t
	}
	if attrs.TimeUpdated != nil {
		t := attrs.TimeUpdated.Time
		timeUpdated = &t
	}

	return &domain.CacheCluster{
		ID:                         val(attrs.ID),
		DisplayName:                val(attrs.DisplayName),
		CompartmentOCID:            val(attrs.CompartmentOCID),
		LifecycleState:             attrs.LifecycleState,
		LifecycleDetails:           val(attrs.LifecycleDetails),
		TimeCreated:                timeCreated,
		TimeUpdated:                timeUpdated,
		NodeCount:                  intVal(attrs.NodeCount),
		NodeMemoryInGBs:            float32Val(attrs.NodeMemoryInGBs),
		SoftwareVersion:            attrs.SoftwareVersion,
		ClusterMode:                attrs.ClusterMode,
		ShardCount:                 intVal(attrs.ShardCount),
		ConfigSetId:                val(attrs.ConfigSetId),
		SubnetId:                   val(attrs.SubnetId),
		NsgIds:                     attrs.NsgIds,
		PrimaryFqdn:                val(attrs.PrimaryFqdn),
		PrimaryEndpointIpAddress:   val(attrs.PrimaryEndpointIpAddress),
		ReplicasFqdn:               val(attrs.ReplicasFqdn),
		ReplicasEndpointIpAddress:  val(attrs.ReplicasEndpointIpAddress),
		DiscoveryFqdn:              val(attrs.DiscoveryFqdn),
		DiscoveryEndpointIpAddress: val(attrs.DiscoveryEndpointIpAddress),
		Nodes:                      attrs.Nodes,
		FreeformTags:               attrs.FreeformTags,
		DefinedTags:                attrs.DefinedTags,
		SystemTags:                 attrs.SystemTags,
	}
}
