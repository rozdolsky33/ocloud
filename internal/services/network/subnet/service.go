package subnet

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/network/subnet"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for subnet operations.
type Service struct {
	subnetRepo    subnet.SubnetRepository
	logger        logr.Logger
	compartmentID string
}

// NewService creates a Service configured with the provided subnet repository, logger, and compartment ID.
func NewService(repo subnet.SubnetRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		subnetRepo:    repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// List retrieves a paginated list of subnets.
func (s *Service) List(ctx context.Context, limit int, pageNum int) ([]subnet.Subnet, int, string, error) {
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
		return []subnet.Subnet{}, totalCount, "", nil
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
func (s *Service) Find(ctx context.Context, namePattern string) ([]subnet.Subnet, error) {
	s.logger.V(logger.Debug).Info("finding subnet with fuzzy search", "pattern", namePattern)

	allSubnets, err := s.subnetRepo.ListSubnets(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all subnets for search: %w", err)
	}

	index, err := util.BuildIndex(allSubnets, func(s subnet.Subnet) any {
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
	var matchedSubnets []subnet.Subnet
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allSubnets) {
			matchedSubnets = append(matchedSubnets, allSubnets[idx])
		}
	}

	s.logger.Info("found subnet", "count", len(matchedSubnets))
	return matchedSubnets, nil
}

// mapToIndexableSubnets converts a subnet.Subnet into a plain struct used for indexing.
// The returned struct has Name set to the subnet's DisplayName lowercased and CIDR set to the subnet's CidrBlock.
func mapToIndexableSubnets(s subnet.Subnet) any {
	return struct {
		Name string
		CIDR string
	}{
		Name: strings.ToLower(s.DisplayName),
		CIDR: s.CidrBlock,
	}
}
