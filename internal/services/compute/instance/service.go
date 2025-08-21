package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service, for instance, operations.
type Service struct {
	instanceRepo  domain.InstanceRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance.
func NewService(repo domain.InstanceRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		instanceRepo:  repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// List retrieves a paginated list of instances.
func (s *Service) List(ctx context.Context, limit int, pageNum int, showImageDetails bool) ([]Instance, int, string, error) {
	s.logger.V(logger.Debug).Info("listing instances", "limit", limit, "pageNum", pageNum)

	allInstances, err := s.instanceRepo.ListInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing instances from repository: %w", err)
	}

	// Manual pagination.
	totalCount := len(allInstances)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []Instance{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allInstances[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	s.logger.Info("completed instance listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search for instances.
func (s *Service) Find(ctx context.Context, searchPattern string, showImageDetails bool) ([]Instance, error) {
	s.logger.V(logger.Debug).Info("finding instances with fuzzy search", "pattern", searchPattern)

	allInstances, err := s.instanceRepo.ListInstances(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all instances for search: %w", err)
	}

	index, err := util.BuildIndex(allInstances, func(inst Instance) any {
		return mapToIndexableInstance(inst)
	})
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}
	s.logger.V(logger.Debug).Info("Search index built successfully.", "numEntries", len(allInstances))

	fields := []string{"Name", "PrimaryIP", "ImageName", "ImageOS"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	s.logger.V(logger.Debug).Info("Fuzzy search completed.", "numMatches", len(matchedIdxs))
	var results []Instance
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allInstances) {
			results = append(results, allInstances[idx])
		}
	}

	s.logger.Info("instance search complete", "matches", len(results))
	return results, nil
}

// mapToIndexableInstance converts a domain.Instance to a struct suitable for indexing.
func mapToIndexableInstance(inst domain.Instance) any {
	return struct {
		Name      string
		PrimaryIP string
		ImageName string
		ImageOS   string
	}{
		Name:      strings.ToLower(inst.DisplayName),
		PrimaryIP: inst.PrimaryIP,
		ImageName: strings.ToLower(inst.ImageName),
		ImageOS:   strings.ToLower(inst.ImageOS),
	}
}
