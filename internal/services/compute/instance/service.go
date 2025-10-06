package instance

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service, for instance, operations.
type Service struct {
	instanceRepo  compute.InstanceRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance.
func NewService(repo compute.InstanceRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		instanceRepo:  repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// ListInstances retrieves a list of instances.
func (s *Service) ListInstances(ctx context.Context) ([]Instance, error) {
	s.logger.V(logger.Debug).Info("listing instances")
	instances, err := s.instanceRepo.ListInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing instances from repository: %w", err)
	}
	return instances, nil
}

// FetchPaginatedInstances retrieves a paginated list of instances.
func (s *Service) FetchPaginatedInstances(ctx context.Context, limit int, pageNum int) ([]Instance, int, string, error) {
	s.logger.V(logger.Debug).Info("listing instances", "limit", limit, "pageNum", pageNum)

	allInstances, err := s.instanceRepo.ListEnrichedInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing instances from repository: %w", err)
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allInstances, limit, pageNum)

	s.logger.Info("completed instance listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// FuzzySearch performs a fuzzy search for instances.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]Instance, error) {
	s.logger.V(logger.Debug).Info("finding instances", "pattern", searchPattern)

	all, err := s.instanceRepo.ListEnrichedInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching instances: %w", err)
	}

	searchableInstances := ToSearchableInstances(all)
	indexMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(searchableInstances, indexMapping)
	if err != nil {
		return nil, fmt.Errorf("build index: %w", err)
	}

	s.logger.V(logger.Debug).Info("index ready", "count", len(all))

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	results := make([]Instance, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(all) {
			results = append(results, all[i])
		}
	}

	return results, nil
}
