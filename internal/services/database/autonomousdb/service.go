package autonomousdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service provides operations and functionalities related to database management, logging, and compartment handling.
type Service struct {
	repo          database.AutonomousDatabaseRepository
	logger        logr.Logger
	compartmentID string
}

// NewService creates a Service configured with the provided repository and application context.
// The returned Service uses the application context's Logger and CompartmentID.
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

	// First try enriched list
	allDatabases, err := s.repo.ListEnrichedAutonomousDatabase(ctx, s.compartmentID)
	if err != nil {
		// Fallback to base list on error (as tests expect)
		allDatabases, err = s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
		if err != nil {
			return nil, 0, "", fmt.Errorf("failed to list autonomous databases: %w", err)
		}
	}

	// If enriched list is empty, also fallback to base list to confirm (as tests expect)
	if len(allDatabases) == 0 {
		var baseErr error
		allDatabases, baseErr = s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
		if baseErr != nil {
			return nil, 0, "", fmt.Errorf("failed to list autonomous databases: %w", baseErr)
		}
	}

	totalCount := len(allDatabases)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []AutonomousDatabase{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allDatabases[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	logger.LogWithLevel(s.logger, logger.Info, "completed database listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search to find autonomous databases matching the given search pattern in their Name field.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]AutonomousDatabase, error) {
	logger.LogWithLevel(s.logger, logger.Trace, "finding database with bleve fuzzy search", "pattern", searchPattern)

	// 1: Fetch all databases
	allDatabases, err := s.repo.ListEnrichedAutonomousDatabase(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all databases: %w", err)
	}
	// 2: Build index
	index, err := util.BuildIndex(allDatabases, func(db database.AutonomousDatabase) any {
		return mapToIndexableDatabase(db)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build index: %w", err)
	}

	// Step 3: Fuzzy search on multiple fields
	fields := []string{"Name"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("failed to fuzzy search index: %w", err)
	}

	var results []AutonomousDatabase
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allDatabases) {
			results = append(results, allDatabases[idx]) // TODO: REDO
		}
	}
	return results, nil
}

// mapToIndexableDatabase converts a database.AutonomousDatabase into an IndexableAutonomousDatabase.
// The returned value contains only the `Name` field, suitable for building search/index structures.
func mapToIndexableDatabase(db database.AutonomousDatabase) IndexableAutonomousDatabase {
	return IndexableAutonomousDatabase{
		Name: db.Name,
	}
}
