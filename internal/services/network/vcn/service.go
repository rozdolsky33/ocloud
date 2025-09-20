package vcn

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/logger"
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

	totalCount := len(allVcn)

	if pageNum <= 0 {
		pageNum = 1
	}

	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []VCN{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allVcn[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

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

//// Find performs a fuzzy search for vcns.
//func (s *Service) Find(ctx context.Context, searchPattern string) ([]*domain.VCN, error) {
//	all, err := s.vcnRepo.ListVcns(ctx, s.compartmentID)
//	if err != nil {
//		return nil, fmt.Errorf("fetching all vcns for search: %w", err)
//	}
//
//	index, err := util.BuildIndex(all, func(vcn *domain.VCN) any {
//		return mapToIndexableVCN(vcn)
//	})
//	if err != nil {
//		return nil, fmt.Errorf("building search index: %w", err)
//	}
//
//	fields := []string{"Name", "OCID"}
//	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
//	if err != nil {
//		return nil, fmt.Errorf("performing fuzzy search: %w", err)
//	}
//	var results []*domain.VCN
//	for _, idx := range matchedIdxs {
//		if idx >= 0 && idx < len(all) {
//			results = append(results, all[idx])
//		}
//	}
//
//	return results, nil
//}

//func ToDTO(v *domain.VCN) *VCNDTO {
//	dto := &VCNDTO{}
//	if v == nil {
//		return dto
//	}
//	dto.OCID = v.OCID
//	dto.DisplayName = v.DisplayName
//	dto.LifecycleState = v.LifecycleState
//	dto.CompartmentID = v.CompartmentID
//	dto.DnsLabel = v.DnsLabel
//	dto.DomainName = v.DomainName
//	dto.CidrBlocks = v.CidrBlocks
//	dto.Ipv6Enabled = v.Ipv6Enabled
//	dto.DhcpOptionsID = v.DhcpOptionsID
//	dto.TimeCreated = v.TimeCreated.UTC().Format(time.RFC3339)
//	return dto
//}

// mapToIndexableVCN converts a domain.VCN to a struct suitable for indexing.
func mapToIndexableVCN(v *VCN) any {
	return struct {
		Name string
		OCID string
	}{
		Name: strings.ToLower(v.DisplayName),
		OCID: strings.ToLower(v.OCID),
	}
}
