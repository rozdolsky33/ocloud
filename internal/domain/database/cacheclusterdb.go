package database

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/v65/redis"
)

// CacheCluster represents an OCI Cache (Redis/Valkey) cluster with its attributes and configuration.
type CacheCluster struct {
	// Identity & lifecycle
	ID               string
	DisplayName      string
	CompartmentOCID  string
	LifecycleState   string
	LifecycleDetails string
	TimeCreated      *time.Time
	TimeUpdated      *time.Time

	// Cluster configuration
	NodeCount       int
	NodeMemoryInGBs float32
	SoftwareVersion string
	ClusterMode     string
	ShardCount      int
	ConfigSetId     string

	// Networking
	SubnetId                   string
	SubnetName                 string
	VcnID                      string
	VcnName                    string
	NsgIds                     []string
	NsgNames                   []string
	PrimaryFqdn                string
	PrimaryEndpointIpAddress   string
	ReplicasFqdn               string
	ReplicasEndpointIpAddress  string
	DiscoveryFqdn              string
	DiscoveryEndpointIpAddress string

	// Nodes
	Nodes []redis.Node

	// Tags
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
	SystemTags   map[string]map[string]interface{}
}

// CacheClusterRepository defines the interface for interacting with OCI Cache Cluster data.
type CacheClusterRepository interface {
	GetCacheCluster(ctx context.Context, clusterId string) (*CacheCluster, error)
	ListCacheClusters(ctx context.Context, compartmentID string) ([]CacheCluster, error)
	ListEnrichedCacheClusters(ctx context.Context, compartmentID string) ([]CacheCluster, error)
}
