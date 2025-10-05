package oke

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for OKE operations.
type Service struct {
	clusterRepo   compute.ClusterRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance.
func NewService(repo compute.ClusterRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		clusterRepo:   repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}
func (s *Service) ListClusters(ctx context.Context) ([]Cluster, error) {
	s.logger.V(logger.Debug).Info("listing clusters")
	clusters, err := s.clusterRepo.ListClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing clusters from repository: %w", err)
	}
	return clusters, nil
}

// FetchPaginatedClusters retrieves a paginated list of clusters.
func (s *Service) FetchPaginatedClusters(ctx context.Context, limit, pageNum int) ([]Cluster, int, string, error) {
	s.logger.V(logger.Debug).Info("listing clusters", "limit", limit, "pageNum", pageNum)

	allClusters, err := s.clusterRepo.ListClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing clusters from repository: %w", err)
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allClusters, limit, pageNum)

	s.logger.Info("completed cluster listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a case-insensitive search for clusters.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Cluster, error) {
	s.logger.V(logger.Debug).Info("finding clusters with search", "pattern", searchPattern)

	allClusters, err := s.clusterRepo.ListClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all clusters for search: %w", err)
	}

	if searchPattern == "" {
		s.logger.V(logger.Debug).Info("Empty search pattern, returning all clusters.")
		return allClusters, nil
	}

	var matchedClusters []Cluster
	searchPattern = strings.ToLower(searchPattern)

	s.logger.V(logger.Trace).Info("Starting cluster iteration for search.", "totalClusters", len(allClusters))

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
