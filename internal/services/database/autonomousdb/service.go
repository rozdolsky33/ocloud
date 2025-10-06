package autonomousdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service provides operations and functionalities related to database management, logging, and compartment handling.
type Service struct {
	repo          database.AutonomousDatabaseRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance with the provided application context.
func NewService(repo database.AutonomousDatabaseRepository, appCtx *app.ApplicationContext) *Service {
	return &Service{
		repo:          repo,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}
}

// ListAutonomousDb retrieves and returns all databases from the given compartment in the OCI account.
func (s *Service) ListAutonomousDb(ctx context.Context) ([]AutonomousDatabase, error) {
	s.logger.V(logger.Debug).Info("listing autonomous databases")
	databases, err := s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list autonomous databases: %w", err)
	}
	return databases, nil
}

// FetchPaginatedAutonomousDb retrieves a paginated list of databases with given limit and page number parameters.
// It returns the slice of databases, total count, next page token, and an error if encountered.
func (s *Service) FetchPaginatedAutonomousDb(ctx context.Context, limit, pageNum int) ([]AutonomousDatabase, int, string, error) {
	s.logger.V(logger.Debug).Info("listing autonomous databases", "limit", limit, "pageNum", pageNum)

	allDatabases, err := s.repo.ListEnrichedAutonomousDatabase(ctx, s.compartmentID)
	if err != nil {
		allDatabases, err = s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
		if err != nil {
			return nil, 0, "", fmt.Errorf("failed to list autonomous databases: %w", err)
		}
	}

	if len(allDatabases) == 0 {
		var baseErr error
		allDatabases, baseErr = s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
		if baseErr != nil {
			return nil, 0, "", fmt.Errorf("failed to list autonomous databases: %w", baseErr)
		}
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allDatabases, limit, pageNum)

	logger.LogWithLevel(s.logger, logger.Info, "completed database listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// FuzzySearch performs a fuzzy search for Autonomous Databases using the generic search engine.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]AutonomousDatabase, error) {
	logger.LogWithLevel(s.logger, logger.Trace, "finding databases with search", "pattern", searchPattern)
	allDatabases, err := s.repo.ListEnrichedAutonomousDatabase(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all databases: %w", err)
	}
	p := strings.TrimSpace(searchPattern)
	if p == "" {
		return allDatabases, nil
	}

	// Build index using SearchableAutonomousDatabase
	indexables := ToSearchableAutonomousDBs(allDatabases)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	hits, err := search.FuzzySearch(idx, strings.ToLower(p), GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("executing search: %w", err)
	}

	results := make([]AutonomousDatabase, 0, len(hits))
	for _, i := range hits {
		if i >= 0 && i < len(allDatabases) {
			results = append(results, allDatabases[i])
		}
	}

	return results, nil
}
