package compartment

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for compartment operations.
// It depends on the domain repository for data access.
type Service struct {
	compartmentRepo domain.CompartmentRepository
	logger          logr.Logger
	tenancyID       string
}

// NewService initializes and returns a new Service instance.
// It injects the domain repository, decoupling the service from the infrastructure layer.
func NewService(repo domain.CompartmentRepository, logger logr.Logger, tenancyID string) *Service {
	return &Service{
		compartmentRepo: repo,
		logger:          logger,
		tenancyID:       tenancyID,
	}
}

// List retrieves a paginated list of compartments.
// Note: The pagination logic here is simplified. A real implementation might need more robust cursor handling.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]domain.Compartment, int, string, error) {
	s.logger.V(logger.Debug).Info("listing compartments", "limit", limit, "pageNum", pageNum)

	// Fetch all compartments from the repository.
	// The underlying adapter handles the complexity of OCI pagination.
	allCompartments, err := s.compartmentRepo.ListCompartments(ctx, s.tenancyID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing compartments from repository: %w", err)
	}

	// Manual pagination over the full list.
	totalCount := len(allCompartments)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []domain.Compartment{}, totalCount, "", nil // Page number is out of bounds
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allCompartments[start:end]

	// Determine if there is a next page.
	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	s.logger.Info("completed compartment listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search for compartments based on the provided searchPattern.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]domain.Compartment, error) {
	s.logger.V(logger.Debug).Info("finding compartments with fuzzy search", "pattern", searchPattern)

	// Step 1: Fetch all compartments from the repository.
	allCompartments, err := s.compartmentRepo.ListCompartments(ctx, s.tenancyID)
	if err != nil {
		return nil, fmt.Errorf("fetching all compartments for search: %w", err)
	}

	// Step 2: Build the search index from the domain models.
	index, err := util.BuildIndex(allCompartments, func(c domain.Compartment) any {
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
	var results []domain.Compartment
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allCompartments) {
			results = append(results, allCompartments[idx])
		}
	}

	s.logger.Info("compartment search complete", "matches", len(results))
	return results, nil
}

// mapToIndexableCompartment converts a domain.Compartment to a struct suitable for indexing.
func mapToIndexableCompartment(compartment domain.Compartment) any {
	return struct {
		Name        string
		Description string
	}{
		Name:        strings.ToLower(compartment.DisplayName),
		Description: strings.ToLower(compartment.Description),
	}
}
