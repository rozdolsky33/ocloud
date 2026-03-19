package dynamicgroup

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for dynamic group operations.
type Service struct {
	dynamicGroupRepo identity.DynamicGroupRepository
	logger           logr.Logger
	compartmentID    string
}

// NewService initializes and returns a new Service instance.
func NewService(repo identity.DynamicGroupRepository, logger logr.Logger, ocid string) *Service {
	return &Service{
		dynamicGroupRepo: repo,
		logger:           logger,
		compartmentID:    ocid,
	}
}

// FetchPaginateDynamicGroups fetches a page of dynamic groups from the repository.
func (s *Service) FetchPaginateDynamicGroups(ctx context.Context, limit, pageNum int) ([]DynamicGroup, int, string, error) {
	s.logger.V(logger.Debug).Info("listing dynamic groups", "limit", limit, "pageNum", pageNum)

	allDynamicGroups, err := s.dynamicGroupRepo.ListDynamicGroups(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing dynamic groups from repository: %w", err)
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allDynamicGroups, limit, pageNum)

	s.logger.Info("completed dynamic group listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// FuzzySearch performs a fuzzy search for dynamic groups based on the provided searchPattern.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]DynamicGroup, error) {
	s.logger.V(logger.Debug).Info("finding dynamic groups with fuzzy search", "pattern", searchPattern)

	allDynamicGroups, err := s.dynamicGroupRepo.ListDynamicGroups(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all dynamic groups for search: %w", err)
	}

	// Build the search index using the common search package and the dynamic group searcher adapter.
	indexables := ToSearchableDynamicGroups(allDynamicGroups)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}

	results := make([]DynamicGroup, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(allDynamicGroups) {
			results = append(results, allDynamicGroups[i])
		}
	}

	return results, nil
}
