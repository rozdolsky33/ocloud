package policy

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

type Service struct {
	policyRepo    identity.PolicyRepository
	logger        logr.Logger
	CompartmentID string
}

// NewService initializes a new Service instance with the provided application context.
func NewService(repo identity.PolicyRepository, logger logr.Logger, ocid string) *Service {
	return &Service{
		policyRepo:    repo,
		logger:        logger,
		CompartmentID: ocid,
	}
}

func (s *Service) FetchPaginatedPolies(ctx context.Context, limit, pageNum int) ([]Policy, int, string, error) {
	s.logger.V(logger.Debug).Info("listing policies", "limit", limit, "pageNum", pageNum)

	allPolicies, err := s.policyRepo.ListPolicies(ctx, s.CompartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing policies from repository: %w", err)
	}
	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allPolicies, limit, pageNum)
	s.logger.V(logger.Debug).Info("completed policy listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

func (s *Service) ListPolicies(ctx context.Context) ([]identity.Policy, error) {
	s.logger.V(logger.Debug).Info("listing policies")
	policies, err := s.policyRepo.ListPolicies(ctx, s.CompartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing policies from repository: %w", err)
	}
	return policies, nil
}

// FuzzySearch performs a fuzzy search for policies based on the provided search pattern and returns matching policies.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]identity.Policy, error) {
	s.logger.V(logger.Debug).Info("finding policies with fuzzy search", "pattern", searchPattern)

	allPolicies, err := s.policyRepo.ListPolicies(ctx, s.CompartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all policies for search: %w", err)
	}

	// Build the search index using the common search package and the policy searcher adapter.
	indexables := ToSearchablePolicies(allPolicies)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	s.logger.V(logger.Debug).Info("search index built successfully", "numEntries", len(allPolicies))

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	s.logger.V(logger.Debug).Info("fuzzy search completed", "numMatches", len(matchedIdxs))

	results := make([]identity.Policy, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(allPolicies) {
			results = append(results, allPolicies[i])
		}
	}

	return results, nil
}
