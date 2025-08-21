package oke

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// Service is the application-layer service for OKE operations.
type Service struct {
	clusterRepo   domain.ClusterRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance.
func NewService(repo domain.ClusterRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		clusterRepo:   repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// List retrieves a paginated list of clusters.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Cluster, int, string, error) {
	s.logger.V(1).Info("listing clusters", "limit", limit, "pageNum", pageNum)

	allClusters, err := s.clusterRepo.ListClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing clusters from repository: %w", err)
	}

	// Manual pagination.
	totalCount := len(allClusters)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []Cluster{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allClusters[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	s.logger.Info("completed cluster listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a case-insensitive search for clusters.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Cluster, error) {
	s.logger.V(1).Info("finding clusters with search", "pattern", searchPattern)

	allClusters, err := s.clusterRepo.ListClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all clusters for search: %w", err)
	}

	if searchPattern == "" {
		return allClusters, nil
	}

	var matchedClusters []Cluster
	searchPattern = strings.ToLower(searchPattern)

	for _, cluster := range allClusters {
		if strings.Contains(strings.ToLower(cluster.DisplayName), searchPattern) {
			matchedClusters = append(matchedClusters, cluster)
			continue
		}
		for _, nodePool := range cluster.NodePools {
			if strings.Contains(strings.ToLower(nodePool.DisplayName), searchPattern) {
				matchedClusters = append(matchedClusters, cluster)
				break
			}
		}
	}

	s.logger.Info("cluster search complete", "matches", len(matchedClusters))
	return matchedClusters, nil
}
