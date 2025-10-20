package cacheclusterdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service provides operations and functionalities related to HeatWave cache cluster management, logging, and compartment handling.
type Service struct {
	repo          database.CacheClusterRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance with the provided application context.
func NewService(repo database.CacheClusterRepository, appCtx *app.ApplicationContext) *Service {
	return &Service{
		repo:          repo,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}
}

// ListCacheClusters retrieves and returns all HeatWave cache clusters from the given compartment in the OCI account.
func (s *Service) ListCacheClusters(ctx context.Context) ([]CacheCluster, error) {
	s.logger.V(logger.Debug).Info("listing HeatWave cache clusters")
	clusters, err := s.repo.ListCacheClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cache clusters: %w", err)
	}
	return clusters, nil
}

// FetchPaginatedCacheClusters retrieves a paginated list of HeatWave cache clusters with given limit and page number parameters.
// It returns the slice of clusters, total count, next page token, and an error if encountered.
func (s *Service) FetchPaginatedCacheClusters(ctx context.Context, limit, pageNum int) ([]CacheCluster, int, string, error) {
	s.logger.V(logger.Debug).Info("listing HeatWave cache clusters", "limit", limit, "pageNum", pageNum)

	allClusters, err := s.repo.ListEnrichedCacheClusters(ctx, s.compartmentID)
	if err != nil {
		allClusters, err = s.repo.ListCacheClusters(ctx, s.compartmentID)
		if err != nil {
			return nil, 0, "", fmt.Errorf("failed to list cache clusters: %w", err)
		}
	}

	if len(allClusters) == 0 {
		var baseErr error
		allClusters, baseErr = s.repo.ListCacheClusters(ctx, s.compartmentID)
		if baseErr != nil {
			return nil, 0, "", fmt.Errorf("failed to list cache clusters: %w", baseErr)
		}
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allClusters, limit, pageNum)

	logger.LogWithLevel(s.logger, logger.Info, "completed cache cluster listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// FuzzySearch performs a fuzzy search across HeatWave cache clusters using a given search pattern.
// It indexes all searchable cluster fields and returns matching clusters.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]CacheCluster, error) {
	logger.LogWithLevel(s.logger, logger.Trace, "finding cache clusters with search", "pattern", searchPattern)
	allClusters, err := s.repo.ListEnrichedCacheClusters(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all cache clusters: %w", err)
	}
	p := strings.TrimSpace(searchPattern)
	if p == "" {
		return allClusters, nil
	}

	// Build index using SearchableCacheCluster
	indexables := ToSearchableCacheClusters(allClusters)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	hits, err := search.FuzzySearch(idx, strings.ToLower(p), GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("executing search: %w", err)
	}

	results := make([]CacheCluster, 0, len(hits))
	for _, i := range hits {
		if i >= 0 && i < len(allClusters) {
			results = append(results, allClusters[i])
		}
	}

	logger.LogWithLevel(s.logger, logger.Debug, "completed search", "pattern", searchPattern, "totalClusters", len(allClusters), "matchedClusters", len(results))
	return results, nil
}
