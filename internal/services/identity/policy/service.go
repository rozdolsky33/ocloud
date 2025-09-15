package policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

type Service struct {
	policyRepo    domain.PolicyRepository
	logger        logr.Logger
	CompartmentID string
}

// NewService initializes a new Service instance with the provided application context.
func NewService(repo domain.PolicyRepository, logger logr.Logger, CompartmentID string) *Service {
	return &Service{
		policyRepo:    repo,
		logger:        logger,
		CompartmentID: CompartmentID,
	}
}

func (s *Service) FetchPaginatedPolies(ctx context.Context, limit, pageNum int) ([]domain.Policy, int, string, error) {
	s.logger.V(logger.Debug).Info("listing policies", "limit", limit, "pageNum", pageNum)

	allPolicies, err := s.policyRepo.ListPolicies(ctx, s.CompartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing policies from repository: %w", err)
	}
	totalCount := len(allPolicies)
	start := (pageNum - 1) * limit
	end := start + limit
	if start >= totalCount {
		return []domain.Policy{}, totalCount, "", nil
	}
	if end > totalCount {
		end = totalCount
	}
	pagedResults := allPolicies[start:end]
	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}
	s.logger.V(logger.Debug).Info("completed policy listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

func (s *Service) ListPolicies(ctx context.Context) ([]domain.Policy, error) {
	s.logger.V(logger.Debug).Info("listing policies")
	policies, err := s.policyRepo.ListPolicies(ctx, s.CompartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing policies from repository: %w", err)
	}
	return policies, nil
}

// Find performs a fuzzy search for policies based on the provided searchPattern and returns matching policy.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]domain.Policy, error) {
	logger.LogWithLevel(s.logger, logger.Debug, "Finding Policies", "pattern", searchPattern)
	var allPolicies []domain.Policy

	// 1. Fetch all policies in the compartment
	allPolicies, err := s.policyRepo.ListPolicies(ctx, s.CompartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all policies: %w", err)
	}

	// 2. Build index
	index, err := util.BuildIndex(allPolicies, func(p domain.Policy) any {
		return mapToIndexablePolicy(p)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build index: %w", err)
	}

	// 3. Fuzzy search on multiple files
	fields := []string{"Name", "Description", "Statement", "Tags", "TagValues"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("failed to fuzzy search index: %w", err)
	}

	// Return matched policies
	var matchedPolicies []domain.Policy
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allPolicies) {
			matchedPolicies = append(matchedPolicies, allPolicies[idx])
		}
	}

	logger.LogWithLevel(s.logger, logger.Trace, "Found policies", "count", len(matchedPolicies))
	return matchedPolicies, nil
}

// mapToIndexablePolicy transforms a Policy object into an IndexablePolicy object with indexed and searchable fields.
func mapToIndexablePolicy(p domain.Policy) IndexablePolicy {
	flattenedTags, _ := util.FlattenTags(p.FreeformTags, p.DefinedTags)
	tagValues, _ := util.ExtractTagValues(p.FreeformTags, p.DefinedTags)
	joinedStatements := strings.Join(p.Statement, " ") // Concatenate all the statements into one string for fuzzy search/indexing.
	return IndexablePolicy{
		Name:        p.Name,
		Description: p.Description,
		Statement:   joinedStatements,
		Tags:        flattenedTags,
		TagValues:   tagValues,
	}

}
