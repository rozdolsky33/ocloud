package compartment

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for compartment operations.
// It depends on the domain repository for data access.
type Service struct {
	compartmentRepo identity.CompartmentRepository
	logger          logr.Logger
	compartmentID   string
}

// NewService initializes and returns a new Service instance.
// It injects the domain repository, decoupling the service from the infrastructure layer.
func NewService(repo identity.CompartmentRepository, logger logr.Logger, ocid string) *Service {
	return &Service{
		compartmentRepo: repo,
		logger:          logger,
		compartmentID:   ocid,
	}
}

// FetchPaginateCompartments fetches a page of compartments from the repository.
func (s *Service) FetchPaginateCompartments(ctx context.Context, limit, pageNum int) ([]Compartment, int, string, error) {
	s.logger.V(logger.Debug).Info("listing compartments", "limit", limit, "pageNum", pageNum)

	allCompartments, err := s.compartmentRepo.ListCompartments(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing compartments from repository: %w", err)
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allCompartments, limit, pageNum)

	s.logger.Info("completed compartment listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// FuzzySearch performs a fuzzy search for compartments based on the provided searchPattern.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]Compartment, error) {
	s.logger.V(logger.Debug).Info("finding compartments with fuzzy search", "pattern", searchPattern)

	allCompartments, err := s.compartmentRepo.ListCompartments(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all compartments for search: %w", err)
	}

	// Build the search index using the common search package and the compartment searcher adapter.
	indexables := ToSearchableCompartments(allCompartments)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	logger.Logger.V(logger.Debug).Info("Search index built successfully.", "numEntries", len(allCompartments))

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("Fuzzy search completed.", "numMatches", len(matchedIdxs))

	results := make([]Compartment, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(allCompartments) {
			results = append(results, allCompartments[i])
		}
	}

	return results, nil
}
