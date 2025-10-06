package oke

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
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

// FuzzySearch performs a fuzzy search for clusters using the generic search engine.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]Cluster, error) {
	s.logger.V(logger.Debug).Info("finding clusters with search", "pattern", searchPattern)

	allClusters, err := s.clusterRepo.ListClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all clusters for search: %w", err)
	}

	// If no search pattern provided, return all clusters
	p := strings.TrimSpace(searchPattern)
	if p == "" {
		s.logger.V(logger.Debug).Info("empty search pattern, returning all clusters")
		return allClusters, nil
	}

	// Build index using SearchableCluster adapter
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	indexables := ToSearchableClusters(allClusters)
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	// Execute fuzzy search
	hits, err := search.FuzzySearch(idx, p, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("executing search: %w", err)
	}

	if len(hits) == 0 {
		s.logger.V(logger.Debug).Info("no matches found for pattern")
		return nil, nil
	}

	matched := make([]Cluster, 0, len(hits))
	for _, i := range hits {
		if i >= 0 && i < len(allClusters) {
			matched = append(matched, allClusters[i])
		}
	}

	s.logger.Info("cluster search complete", "matches", len(matched))
	return matched, nil
}
