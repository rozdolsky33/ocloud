package subnet

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for subnet operations.
type Service struct {
	subnetRepo    domain.SubnetRepository
	logger        logr.Logger
	compartmentID string
}

// NewService creates and initializes a new Service instance.
func NewService(repo domain.SubnetRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		subnetRepo:    repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// List retrieves a paginated list of subnets.
func (s *Service) List(ctx context.Context, limit int, pageNum int) ([]Subnet, int, string, error) {
	s.logger.V(logger.Debug).Info("listing subnets", "limit", limit, "pageNum", pageNum)

	allSubnets, err := s.subnetRepo.ListSubnets(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing subnets from repository: %w", err)
	}

	// Manual pagination.
	totalCount := len(allSubnets)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []Subnet{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allSubnets[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	s.logger.Info("completed subnet listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find retrieves a slice of subnets whose attributes match the provided name pattern using fuzzy search.
func (s *Service) Find(ctx context.Context, namePattern string) ([]Subnet, error) {
	s.logger.V(logger.Debug).Info("finding subnet with fuzzy search", "pattern", namePattern)

	allSubnets, err := s.subnetRepo.ListSubnets(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all subnets for search: %w", err)
	}

	index, err := util.BuildIndex(allSubnets, func(s Subnet) any {
		return mapToIndexableSubnets(s)
	})
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("Search index built successfully.", "numEntries", len(allSubnets))

	fields := []string{"Name", "CIDR"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(namePattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("Fuzzy search completed.", "numMatches", len(matchedIdxs))
	var matchedSubnets []Subnet
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allSubnets) {
			matchedSubnets = append(matchedSubnets, allSubnets[idx])
		}
	}

	s.logger.Info("found subnet", "count", len(matchedSubnets))
	return matchedSubnets, nil
}

// mapToIndexableSubnets converts a domain.Subnet to a struct suitable for indexing.
func mapToIndexableSubnets(s domain.Subnet) any {
	return struct {
		Name string
		CIDR string
	}{
		Name: strings.ToLower(s.DisplayName),
		CIDR: s.CIDRBlock,
	}
}
