package mapping

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/redis"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheClusterAttributesFromOCIRedisCluster(t *testing.T) {
	// Setup test data
	id := "ocid1.rediscluster.oc1.iad.test123"
	displayName := "test-cluster"
	compartmentId := "ocid1.compartment.oc1..test"
	nodeCount := 3
	nodeMemory := float32(16.0)
	primaryFqdn := "test.redis.us-ashburn-1.oci.oraclecloud.com"
	primaryIP := "10.0.1.100"
	replicasFqdn := "test-replicas.redis.us-ashburn-1.oci.oraclecloud.com"
	replicasIP := "10.0.1.101"
	subnetId := "ocid1.subnet.oc1.iad.subnet123"
	lifecycleDetails := "Cluster is active"
	configSetId := "ocid1.oci-cache-config-set.oc1.iad.config123"
	shardCount := 2
	nsgIds := []string{"ocid1.networksecuritygroup.oc1.iad.nsg1"}

	now := common.SDKTime{Time: time.Now()}

	nodeDisplayName := "test-cluster-0"
	nodeEndpoint := "test-node-0.redis.us-ashburn-1.oci.oraclecloud.com"
	nodeIP := "10.0.1.200"

	cluster := redis.RedisCluster{
		Id:                        &id,
		DisplayName:               &displayName,
		CompartmentId:             &compartmentId,
		NodeCount:                 &nodeCount,
		NodeMemoryInGBs:           &nodeMemory,
		PrimaryFqdn:               &primaryFqdn,
		PrimaryEndpointIpAddress:  &primaryIP,
		ReplicasFqdn:              &replicasFqdn,
		ReplicasEndpointIpAddress: &replicasIP,
		SubnetId:                  &subnetId,
		LifecycleState:            redis.RedisClusterLifecycleStateActive,
		LifecycleDetails:          &lifecycleDetails,
		TimeCreated:               &now,
		TimeUpdated:               &now,
		SoftwareVersion:           redis.RedisClusterSoftwareVersionValkey72,
		ClusterMode:               redis.RedisClusterClusterModeSharded,
		ShardCount:                &shardCount,
		OciCacheConfigSetId:       &configSetId,
		NsgIds:                    nsgIds,
		NodeCollection: &redis.NodeCollection{
			Items: []redis.Node{
				{
					DisplayName:              &nodeDisplayName,
					PrivateEndpointFqdn:      &nodeEndpoint,
					PrivateEndpointIpAddress: &nodeIP,
				},
			},
		},
		FreeformTags: map[string]string{"env": "test"},
		DefinedTags:  map[string]map[string]interface{}{"namespace": {"key": "value"}},
		SystemTags:   map[string]map[string]interface{}{"system": {"tag": "val"}},
	}

	// Execute
	attrs := NewCacheClusterAttributesFromOCIRedisCluster(cluster)

	// Assert
	assert.NotNil(t, attrs)
	assert.Equal(t, &id, attrs.ID)
	assert.Equal(t, &displayName, attrs.DisplayName)
	assert.Equal(t, &compartmentId, attrs.CompartmentOCID)
	assert.Equal(t, "ACTIVE", attrs.LifecycleState)
	assert.Equal(t, &lifecycleDetails, attrs.LifecycleDetails)
	assert.Equal(t, &now, attrs.TimeCreated)
	assert.Equal(t, &now, attrs.TimeUpdated)
	assert.Equal(t, &nodeCount, attrs.NodeCount)
	assert.Equal(t, &nodeMemory, attrs.NodeMemoryInGBs)
	assert.Equal(t, "VALKEY_7_2", attrs.SoftwareVersion)
	assert.Equal(t, "SHARDED", attrs.ClusterMode)
	assert.Equal(t, &shardCount, attrs.ShardCount)
	assert.Equal(t, &configSetId, attrs.ConfigSetId)
	assert.Equal(t, &subnetId, attrs.SubnetId)
	assert.Equal(t, nsgIds, attrs.NsgIds)
	assert.Equal(t, &primaryFqdn, attrs.PrimaryFqdn)
	assert.Equal(t, &primaryIP, attrs.PrimaryEndpointIpAddress)
	assert.Equal(t, &replicasFqdn, attrs.ReplicasFqdn)
	assert.Equal(t, &replicasIP, attrs.ReplicasEndpointIpAddress)
	assert.Len(t, attrs.Nodes, 1)
	assert.Equal(t, &nodeDisplayName, attrs.Nodes[0].DisplayName)
	assert.Equal(t, map[string]string{"env": "test"}, attrs.FreeformTags)
	assert.NotNil(t, attrs.DefinedTags)
	assert.NotNil(t, attrs.SystemTags)
}

func TestNewCacheClusterAttributesFromOCIRedisClusterSummary(t *testing.T) {
	// Setup test data
	id := "ocid1.rediscluster.oc1.iad.summary123"
	displayName := "summary-cluster"
	compartmentId := "ocid1.compartment.oc1..summary"
	nodeCount := 1
	nodeMemory := float32(8.0)
	primaryFqdn := "summary.redis.us-ashburn-1.oci.oraclecloud.com"
	primaryIP := "10.0.2.100"
	subnetId := "ocid1.subnet.oc1.iad.subnet456"

	now := common.SDKTime{Time: time.Now()}

	summary := redis.RedisClusterSummary{
		Id:                        &id,
		DisplayName:               &displayName,
		CompartmentId:             &compartmentId,
		NodeCount:                 &nodeCount,
		NodeMemoryInGBs:           &nodeMemory,
		PrimaryFqdn:               &primaryFqdn,
		PrimaryEndpointIpAddress:  &primaryIP,
		ReplicasFqdn:              &primaryFqdn,
		ReplicasEndpointIpAddress: &primaryIP,
		SubnetId:                  &subnetId,
		LifecycleState:            redis.RedisClusterLifecycleStateActive,
		TimeCreated:               &now,
		TimeUpdated:               &now,
		SoftwareVersion:           redis.RedisClusterSoftwareVersionRedis70,
		ClusterMode:               redis.RedisClusterClusterModeNonsharded,
		FreeformTags:              map[string]string{"tier": "free"},
	}

	// Execute
	attrs := NewCacheClusterAttributesFromOCIRedisClusterSummary(summary)

	// Assert
	assert.NotNil(t, attrs)
	assert.Equal(t, &id, attrs.ID)
	assert.Equal(t, &displayName, attrs.DisplayName)
	assert.Equal(t, &compartmentId, attrs.CompartmentOCID)
	assert.Equal(t, "ACTIVE", attrs.LifecycleState)
	assert.Equal(t, &nodeCount, attrs.NodeCount)
	assert.Equal(t, &nodeMemory, attrs.NodeMemoryInGBs)
	assert.Equal(t, "REDIS_7_0", attrs.SoftwareVersion)
	assert.Equal(t, "NONSHARDED", attrs.ClusterMode)
	assert.Equal(t, &subnetId, attrs.SubnetId)
	assert.Nil(t, attrs.Nodes, "summary should not have nodes")
	assert.Equal(t, map[string]string{"tier": "free"}, attrs.FreeformTags)
}

func TestNewDomainCacheClusterFromAttrs(t *testing.T) {
	// Setup test data
	id := "ocid1.rediscluster.oc1.iad.domain123"
	displayName := "domain-cluster"
	compartmentId := "ocid1.compartment.oc1..domain"
	lifecycleDetails := "Running normally"
	nodeCount := 3
	nodeMemory := float32(16.0)
	shardCount := 2
	configSetId := "ocid1.oci-cache-config-set.oc1.iad.config456"
	subnetId := "ocid1.subnet.oc1.iad.subnet789"
	primaryFqdn := "domain.redis.us-ashburn-1.oci.oraclecloud.com"
	primaryIP := "10.0.3.100"
	replicasFqdn := "domain-replicas.redis.us-ashburn-1.oci.oraclecloud.com"
	replicasIP := "10.0.3.101"
	discoveryFqdn := "domain-discovery.redis.us-ashburn-1.oci.oraclecloud.com"
	discoveryIP := "10.0.3.102"

	now := common.SDKTime{Time: time.Now()}

	attrs := &CacheClusterAttributes{
		ID:                         &id,
		DisplayName:                &displayName,
		CompartmentOCID:            &compartmentId,
		LifecycleState:             "ACTIVE",
		LifecycleDetails:           &lifecycleDetails,
		TimeCreated:                &now,
		TimeUpdated:                &now,
		NodeCount:                  &nodeCount,
		NodeMemoryInGBs:            &nodeMemory,
		SoftwareVersion:            "VALKEY_7_2",
		ClusterMode:                "SHARDED",
		ShardCount:                 &shardCount,
		ConfigSetId:                &configSetId,
		SubnetId:                   &subnetId,
		NsgIds:                     []string{"nsg1", "nsg2"},
		PrimaryFqdn:                &primaryFqdn,
		PrimaryEndpointIpAddress:   &primaryIP,
		ReplicasFqdn:               &replicasFqdn,
		ReplicasEndpointIpAddress:  &replicasIP,
		DiscoveryFqdn:              &discoveryFqdn,
		DiscoveryEndpointIpAddress: &discoveryIP,
		Nodes:                      []redis.Node{},
		FreeformTags:               map[string]string{"env": "prod"},
		DefinedTags:                map[string]map[string]interface{}{"ns": {"k": "v"}},
		SystemTags:                 map[string]map[string]interface{}{"sys": {"t": "v"}},
	}

	// Execute
	cluster := NewDomainCacheClusterFromAttrs(attrs)

	// Assert
	assert.NotNil(t, cluster)
	assert.Equal(t, id, cluster.ID)
	assert.Equal(t, displayName, cluster.DisplayName)
	assert.Equal(t, compartmentId, cluster.CompartmentOCID)
	assert.Equal(t, "ACTIVE", cluster.LifecycleState)
	assert.Equal(t, lifecycleDetails, cluster.LifecycleDetails)
	assert.NotNil(t, cluster.TimeCreated)
	assert.Equal(t, now.Time, *cluster.TimeCreated)
	assert.NotNil(t, cluster.TimeUpdated)
	assert.Equal(t, now.Time, *cluster.TimeUpdated)
	assert.Equal(t, nodeCount, cluster.NodeCount)
	assert.Equal(t, nodeMemory, cluster.NodeMemoryInGBs)
	assert.Equal(t, "VALKEY_7_2", cluster.SoftwareVersion)
	assert.Equal(t, "SHARDED", cluster.ClusterMode)
	assert.Equal(t, shardCount, cluster.ShardCount)
	assert.Equal(t, configSetId, cluster.ConfigSetId)
	assert.Equal(t, subnetId, cluster.SubnetId)
	assert.Equal(t, []string{"nsg1", "nsg2"}, cluster.NsgIds)
	assert.Equal(t, primaryFqdn, cluster.PrimaryFqdn)
	assert.Equal(t, primaryIP, cluster.PrimaryEndpointIpAddress)
	assert.Equal(t, replicasFqdn, cluster.ReplicasFqdn)
	assert.Equal(t, replicasIP, cluster.ReplicasEndpointIpAddress)
	assert.Equal(t, discoveryFqdn, cluster.DiscoveryFqdn)
	assert.Equal(t, discoveryIP, cluster.DiscoveryEndpointIpAddress)
	assert.Equal(t, map[string]string{"env": "prod"}, cluster.FreeformTags)
	assert.NotNil(t, cluster.DefinedTags)
	assert.NotNil(t, cluster.SystemTags)
}

func TestNewDomainCacheClusterFromAttrs_WithNilValues(t *testing.T) {
	// Test with nil pointers to ensure safe handling
	attrs := &CacheClusterAttributes{
		ID:              nil,
		DisplayName:     nil,
		CompartmentOCID: nil,
		LifecycleState:  "CREATING",
		SoftwareVersion: "REDIS_7_0",
		ClusterMode:     "NONSHARDED",
	}

	// Execute
	cluster := NewDomainCacheClusterFromAttrs(attrs)

	// Assert
	assert.NotNil(t, cluster)
	assert.Equal(t, "", cluster.ID)
	assert.Equal(t, "", cluster.DisplayName)
	assert.Equal(t, "", cluster.CompartmentOCID)
	assert.Equal(t, "CREATING", cluster.LifecycleState)
	assert.Equal(t, 0, cluster.NodeCount)
	assert.Equal(t, float32(0), cluster.NodeMemoryInGBs)
	assert.Equal(t, 0, cluster.ShardCount)
	assert.Nil(t, cluster.TimeCreated)
	assert.Nil(t, cluster.TimeUpdated)
}

func TestNewDomainCacheClusterFromAttrs_WithEmptyNodes(t *testing.T) {
	id := "test-id"
	displayName := "test-cluster"

	attrs := &CacheClusterAttributes{
		ID:              &id,
		DisplayName:     &displayName,
		LifecycleState:  "ACTIVE",
		SoftwareVersion: "VALKEY_7_2",
		ClusterMode:     "NONSHARDED",
		Nodes:           []redis.Node{},
	}

	// Execute
	cluster := NewDomainCacheClusterFromAttrs(attrs)

	// Assert
	assert.NotNil(t, cluster)
	assert.NotNil(t, cluster.Nodes)
	assert.Len(t, cluster.Nodes, 0)
}
