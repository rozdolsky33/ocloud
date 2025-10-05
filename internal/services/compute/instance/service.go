package instance

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/logger"
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
func (s *Service) ListInstances(ctx context.Context) ([]compute.Instance, error) {
	s.logger.V(logger.Debug).Info("listing instances")
	instances, err := s.instanceRepo.ListInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing instances from repository: %w", err)
	}
	return instances, nil
}

// FetchPaginatedInstances retrieves a paginated list of instances.
func (s *Service) FetchPaginatedInstances(ctx context.Context, limit int, pageNum int) ([]compute.Instance, int, string, error) {
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
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]compute.Instance, error) {
	s.logger.V(logger.Debug).Info("finding instances", "pattern", searchPattern)

	all, err := s.instanceRepo.ListEnrichedInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching instances: %w", err)
	}

	idx, err := BuildIndex(all)
	if err != nil {
		return nil, fmt.Errorf("build index: %w", err)
	}
	s.logger.V(logger.Debug).Info("index ready", "count", len(all))

	matchedIdxs, err := FuzzySearchInstances(idx, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	results := make([]compute.Instance, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(all) {
			results = append(results, all[i])
		}
	}
	s.logger.Info("instance search complete", "matches", len(results))
	return results, nil
}
