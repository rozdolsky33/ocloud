package vcn

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for vcn operations.
type Service struct {
	vcnRepo       domain.VCNRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance.
func NewService(repo domain.VCNRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		vcnRepo:       repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// FetchPaginatedVCNs retrieves a paginated list of vcns.
func (s *Service) FetchPaginatedVCNs(ctx context.Context, limit, pageNum int) ([]VCN, int, string, error) {
	s.logger.V(logger.Debug).Info("listing vcns", "limit", limit, "pageNum", pageNum)
	allVcn, err := s.vcnRepo.ListEnrichedVcns(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing vcns from repository: %w", err)
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allVcn, limit, pageNum)

	return pagedResults, totalCount, nextPageToken, nil
}

// ListVcns retrieves a list of vcns.
func (s *Service) ListVcns(ctx context.Context) ([]VCN, error) {
	s.logger.V(logger.Debug).Info("listing vcns")
	allVcn, err := s.vcnRepo.ListEnrichedVcns(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing vcns from repository: %w", err)
	}
	return allVcn, nil
}

// FuzzySearch performs a fuzzy search for vcns.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]VCN, error) {
	all, err := s.vcnRepo.ListEnrichedVcns(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all VCNs for search: %w", err)
	}

	// Build the search index using the common search package and the VCN searcher adapter.
	indexables := ToSearchableVCNs(all)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}

	results := make([]VCN, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(all) {
			results = append(results, all[i])
		}
	}
	return results, nil
}
