package compartment

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
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
// NewService creates a Service configured with the given identity.CompartmentRepository, logger, and compartment OCID.
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

	totalCount := len(allCompartments)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []Compartment{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allCompartments[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	s.logger.Info("completed compartment listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search for compartments based on the provided searchPattern.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Compartment, error) {
	s.logger.V(logger.Debug).Info("finding compartments with fuzzy search", "pattern", searchPattern)

	// Step 1: Fetch all compartments from the repository.
	allCompartments, err := s.compartmentRepo.ListCompartments(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all compartments for search: %w", err)
	}

	// Step 2: Build the search index from the domain models.
	index, err := util.BuildIndex(allCompartments, func(c Compartment) any {
		return mapToIndexableCompartment(c)
	})
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("Search index built successfully.", "numEntries", len(allCompartments))

	// Step 3: Perform the fuzzy search.
	fields := []string{"Name", "Description"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("Fuzzy search completed.", "numMatches", len(matchedIdxs))
	var results []Compartment
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allCompartments) {
			results = append(results, allCompartments[idx])
		}
	}

	s.logger.Info("compartment search complete", "matches", len(results))
	return results, nil
}

// mapToIndexableCompartment converts a domain.Compartment to a struct suitable for indexing.
func mapToIndexableCompartment(compartment Compartment) any {
	return struct {
		Name        string
		Description string
	}{
		Name:        strings.ToLower(compartment.DisplayName),
		Description: strings.ToLower(compartment.Description),
	}
}
