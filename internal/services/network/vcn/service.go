package vcn

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/logger"
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

// Find performs a fuzzy search for vcns.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]VCN, error) {
	all, err := s.vcnRepo.ListEnrichedVcns(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all vcns for search: %w", err)
	}

	index, err := util.BuildIndex(all, func(vcn VCN) any {
		return mapToIndexableVCN(&vcn)
	})
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	fields := []string{"Name", "OCID", "DNSLabel", "DomainName", "CidrBlocks", "TagText", "TagValues"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	var results []VCN
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(all) {
			results = append(results, all[idx])
		}
	}

	return results, nil
}

// mapToIndexableVCN converts a domain.VCN to a struct suitable for indexing.
func mapToIndexableVCN(v *VCN) any {
	// Flatten tags for better search coverage
	tagText, _ := util.FlattenTags(v.FreeformTags, v.DefinedTags)
	tagValues, _ := util.ExtractTagValues(v.FreeformTags, v.DefinedTags)
	return struct {
		Name       string
		OCID       string
		DNSLabel   string
		DomainName string
		TagText    string
		TagValues  string
	}{
		Name:       strings.ToLower(v.DisplayName),
		OCID:       strings.ToLower(v.OCID),
		DNSLabel:   strings.ToLower(v.DnsLabel),
		DomainName: strings.ToLower(v.DomainName),
		TagText:    tagText,
		TagValues:  tagValues,
	}
}
