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

// NewService creates a new Service instance with an OCI container engine client using the provided ApplicationContext.
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

// List retrieves OKE clusters within the specified compartment with pagination support.
// It enriches each cluster with its associated node pools.
// Returns a slice of Cluster objects, total count, next page token, and an error, if any.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Cluster, int, string, error) {
	logger.LogWithLevel(s.logger, 3, "Listing clusters with pagination",
		"limit", limit,
		"pageNum", pageNum)

	var clusters []Cluster
	var nextPageToken string
	var totalCount int

	// Create a request
	request := containerengine.ListClustersRequest{
		CompartmentId: &s.compartmentID,
	}

	// Add limit parameter if specified
	if limit > 0 {
		limitInt := limit
		request.Limit = &limitInt
		logger.LogWithLevel(s.logger, 3, "Setting limit parameter", "limit", limit)
	}

	// If pageNum > 1, we need to fetch the appropriate page token
	if pageNum > 1 && limit > 0 {
		logger.LogWithLevel(s.logger, 3, "Calculating page token for page", "pageNum", pageNum)

		// We need to fetch page tokens until we reach the desired page
		var page *string // Initialize as nil
		currentPage := 1

		for currentPage < pageNum {
			// Fetch just the page token, not actual data
			// Use the same limit to ensure consistent pagination
			tokenRequest := containerengine.ListClustersRequest{
				CompartmentId: &s.compartmentID,
			}

			// Only set the Page field if we have a valid page token
			if page != nil {
				tokenRequest.Page = page
			}
			if limit > 0 {
				limitInt := limit
				tokenRequest.Limit = &limitInt
			}

			resp, err := s.containerEngineClient.ListClusters(ctx, tokenRequest)
			if err != nil {
				return nil, 0, "", fmt.Errorf("fetching page token: %w", err)
			}

			// If there's no next page, we've reached the end
			if resp.OpcNextPage == nil {
				logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return an empty result since the requested page is beyond available data
				return []Cluster{}, 0, "", nil
			}
			// Move to the next page
			page = resp.OpcNextPage
			currentPage++
		}
		// Set the page token for the actual request
		request.Page = page
		logger.LogWithLevel(s.logger, 3, "Using page token for page", "pageNum", pageNum, "token", page)
	}

	// Fetch clusters for the requested page
	response, err := s.containerEngineClient.ListClusters(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing clusters: %w", err)
	}

	// Set the total count to the number of clusters returned
	// If we have a next page, this is an estimate
	totalCount = len(response.Items)
	// If we have a next page, we know there are more clusters
	if response.OpcNextPage != nil {
		// Estimate total count based on current page and items per page
		totalCount = pageNum*limit + limit
	}

	// Save the next page token if available
	if response.OpcNextPage != nil {
		nextPageToken = *response.OpcNextPage
		logger.LogWithLevel(s.logger, 3, "Next page token", "token", nextPageToken)
	}

	for _, cluster := range response.Items {
		// Create a cluster without node pools first
		clusterObj := mapToCluster(cluster)

		// Get node pools for this cluster
		nodePools, err := s.getClusterNodePools(ctx, *cluster.Id)
		if err != nil {
			return nil, 0, "", fmt.Errorf("listing node pools: %w", err)
		}

		// Assign node pools to the cluster
		clusterObj.NodePools = nodePools

		// Add the cluster to the result
		clusters = append(clusters, clusterObj)
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount
	logger.LogWithLevel(s.logger, 2, "Completed cluster listing with pagination",
		"returnedCount", len(clusters),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	return clusters, totalCount, nextPageToken, nil
}

// fetchAllClusters retrieves all clusters within the specified compartment using pagination.
// It returns a slice of Cluster objects and an error, if any.
func (s *Service) fetchAllClusters(ctx context.Context) ([]Cluster, error) {
	logger.LogWithLevel(s.logger, 3, "Fetching all clusters")

	var allClusters []Cluster
	var page *string // Initialize as nil

	for {
		// Create a request with pagination
		request := containerengine.ListClustersRequest{
			CompartmentId: &s.compartmentID,
		}

		// Only set the Page field if we have a valid page token
		if page != nil {
			request.Page = page
		}

		// Fetch clusters for the current page
		response, err := s.containerEngineClient.ListClusters(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("listing clusters: %w", err)
		}

		// Process clusters from this page
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
			allClusters = append(allClusters, clusterObj)
		}

		// If there's no next page, we're done
		if response.OpcNextPage == nil {
			break
		}

		// Move to the next page
		page = response.OpcNextPage
	}

	logger.LogWithLevel(s.logger, 2, "Fetched all clusters", "count", len(allClusters))
	return allClusters, nil
}

// Find searches for OKE clusters matching the given pattern within the compartment.
// It performs a case-insensitive search on cluster names and node pool names.
// If searchPattern is empty, it returns all clusters.
// Returns a slice of matching Cluster objects and an error, if any.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Cluster, error) {
	logger.LogWithLevel(s.logger, 1, "Finding clusters", "pattern", searchPattern)

	// First, get all clusters using fetchAllClusters
	clusters, err := s.fetchAllClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing clusters for search: %w", err)
	}

	// If a search pattern is empty, return all clusters
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

// mapToCluster maps an OCI ClusterSummary to our shared Cluster model.
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

// mapToNodePool maps an OCI NodePoolSummary to our shared NodePool model.
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
