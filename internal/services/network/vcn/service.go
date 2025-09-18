package vcn

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"strings"
	"time"
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

// GetVcn fetches a VCN by ID and prints its summary (or JSON).
func (s *Service) GetVcn(ctx context.Context, vcnID string) (*domain.VCN, error) {
	return s.vcnRepo.GetVcn(ctx, vcnID)
}

// ListVcns lists all VCNs in a compartment.
func (s *Service) ListVcns(ctx context.Context) ([]*domain.VCN, error) {
	return s.vcnRepo.ListVcns(ctx, s.compartmentID)
}

// Find performs a fuzzy search for vcns.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]*domain.VCN, error) {
	all, err := s.vcnRepo.ListVcns(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all vcns for search: %w", err)
	}

	index, err := util.BuildIndex(all, func(vcn *domain.VCN) any {
		return mapToIndexableVCN(vcn)
	})
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	fields := []string{"Name", "OCID"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	var results []*domain.VCN
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(all) {
			results = append(results, all[idx])
		}
	}

	return results, nil
}

func ToDTO(v *domain.VCN) *VCNDTO {
	dto := &VCNDTO{}
	if v == nil {
		return dto
	}
	dto.OCID = v.OCID
	dto.DisplayName = v.DisplayName
	dto.LifecycleState = v.LifecycleState
	dto.CompartmentID = v.CompartmentID
	dto.DnsLabel = v.DnsLabel
	dto.DomainName = v.DomainName
	dto.CidrBlocks = v.CidrBlocks
	dto.Ipv6Enabled = v.Ipv6Enabled
	dto.DhcpOptionsID = v.DhcpOptionsID
	dto.TimeCreated = v.TimeCreated.UTC().Format(time.RFC3339)
	dto.FreeformTags = v.FreeformTags
	return dto
}

// mapToIndexableVCN converts a domain.VCN to a struct suitable for indexing.
func mapToIndexableVCN(v *domain.VCN) any {
	return struct {
		Name string
		OCID string
	}{
		Name: strings.ToLower(v.DisplayName),
		OCID: strings.ToLower(v.OCID),
	}
}
