package oke

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// NewService creates a new Service instance with OCI container engine client using the provided ApplicationContext.
// Returns a Service pointer and an error if the initialization fails.
func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	cfg := appCtx.Provider
	cec, err := oci.NewContainerEngineClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create container engine client: %w", err)
	}
	return &Service{
		containerEngineClient: cec,
		logger:                appCtx.Logger,
		compartmentID:         appCtx.CompartmentID,
	}, nil
}

// List retrieves all OKE clusters within the specified compartment.
// It enriches each cluster with its associated node pools.
// Returns a slice of Cluster objects and an error, if any.
func (s *Service) List(ctx context.Context) ([]Cluster, error) {
	logger.LogWithLevel(s.logger, 3, "Listing clusters")

	var clusters []Cluster

	// Create a request
	request := containerengine.ListClustersRequest{
		CompartmentId: &s.compartmentID,
	}
	response, err := s.containerEngineClient.ListClusters(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("listing clusters: %w", err)
	}

	for _, cluster := range response.Items {
		// Create a cluster without node pools first
		clusterObj := mapToCluster(cluster)

		// Get node pools for this cluster
		nodePools, err := s.getClusterNodePools(ctx, *cluster.Id)
		if err != nil {
			return nil, fmt.Errorf("listing node pools: %w", err)
		}

		// Assign node pools to the cluster
		clusterObj.NodePools = nodePools

		// Add the cluster to the result
		clusters = append(clusters, clusterObj)
	}

	return clusters, nil
}

// Find searches for OKE clusters matching the given pattern within the compartment.
// It performs a case-insensitive search on cluster names and node pool names.
// If searchPattern is empty, it returns all clusters.
// Returns a slice of matching Cluster objects and an error, if any.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Cluster, error) {
	logger.LogWithLevel(s.logger, 1, "Finding clusters", "pattern", searchPattern)

	// First, get all clusters
	clusters, err := s.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing clusters for search: %w", err)
	}

	// If search pattern is empty, return all clusters
	if searchPattern == "" {
		return clusters, nil
	}

	// Filter clusters by name (case-insensitive)
	var matchedClusters []Cluster
	searchPattern = strings.ToLower(searchPattern)

	for _, cluster := range clusters {
		// Check if the cluster name contains the search pattern
		if strings.Contains(strings.ToLower(cluster.Name), searchPattern) {
			matchedClusters = append(matchedClusters, cluster)
			continue
		}

		// Check if any node pool name contains the search pattern
		for _, nodePool := range cluster.NodePools {
			if strings.Contains(strings.ToLower(nodePool.Name), searchPattern) {
				matchedClusters = append(matchedClusters, cluster)
				break
			}
		}
	}

	logger.LogWithLevel(s.logger, 2, "Found clusters", "count", len(matchedClusters))
	return matchedClusters, nil
}

// getClusterNodePools retrieves all node pools associated with the specified cluster.
// It returns a slice of NodePool objects and an error, if any.
func (s *Service) getClusterNodePools(ctx context.Context, clusterID string) ([]NodePool, error) {
	logger.LogWithLevel(s.logger, 3, "Getting node pools for cluster", "clusterID", clusterID)

	var clusterNodePools []NodePool
	nodePools, err := s.containerEngineClient.ListNodePools(ctx, containerengine.ListNodePoolsRequest{
		CompartmentId: common.String(s.compartmentID),
		ClusterId:     common.String(clusterID),
	})
	if err != nil {
		return nil, fmt.Errorf("listing node pools: %w", err)
	}

	for _, nodePool := range nodePools.Items {
		clusterNodePools = append(clusterNodePools, mapToNodePool(nodePool))
	}

	logger.LogWithLevel(s.logger, 3, "Found node pools for cluster", "clusterID", clusterID, "count", len(clusterNodePools))
	return clusterNodePools, nil
}

// mapToCluster maps an OCI ClusterSummary to our internal Cluster model.
// It initializes the NodePools field as an empty slice, which will be populated later.
func mapToCluster(cluster containerengine.ClusterSummary) Cluster {
	return Cluster{
		ID:              *cluster.Id,
		Name:            *cluster.Name,
		CreatedAt:       cluster.Metadata.TimeCreated.String(),
		Version:         *cluster.KubernetesVersion,
		State:           cluster.LifecycleState,
		PrivateEndpoint: *cluster.Endpoints.PrivateEndpoint,
		VcnID:           *cluster.VcnId,
		NodePools:       []NodePool{},
		OKETags: util.ResourceTags{
			FreeformTags: cluster.FreeformTags,
			DefinedTags:  cluster.DefinedTags,
		},
	}
}

// mapToNodePool maps an OCI NodePoolSummary to our internal NodePool model.
func mapToNodePool(nodePool containerengine.NodePoolSummary) NodePool {
	var nodeCount int
	if nodePool.NodeConfigDetails != nil && nodePool.NodeConfigDetails.Size != nil {
		nodeCount = *nodePool.NodeConfigDetails.Size
	} else if nodePool.QuantityPerSubnet != nil {
		nodeCount = *nodePool.QuantityPerSubnet * len(nodePool.SubnetIds)
	}

	// Extract image details from NodeSourceDetails
	image := ""
	if details, ok := nodePool.NodeSourceDetails.(containerengine.NodeSourceViaImageDetails); ok && details.ImageId != nil {
		image = *details.ImageId
	}

	// Optional custom logic for parsing shapeConfig
	ocpus := ""
	memory := ""
	if nodePool.NodeShapeConfig != nil {
		if nodePool.NodeShapeConfig.Ocpus != nil {
			ocpus = fmt.Sprintf("%.1f", *nodePool.NodeShapeConfig.Ocpus)
		}
		if nodePool.NodeShapeConfig.MemoryInGBs != nil {
			memory = fmt.Sprintf("%.0f", *nodePool.NodeShapeConfig.MemoryInGBs)
		}
	}

	return NodePool{
		Name:      *nodePool.Name,
		ID:        *nodePool.Id,
		Version:   *nodePool.KubernetesVersion,
		State:     nodePool.LifecycleState,
		NodeShape: *nodePool.NodeShape,
		NodeCount: nodeCount,
		Image:     image,
		Ocpus:     ocpus,
		MemoryGB:  memory,
		NodeTags: util.ResourceTags{
			FreeformTags: nodePool.FreeformTags,
			DefinedTags:  nodePool.DefinedTags,
		},
	}
}
