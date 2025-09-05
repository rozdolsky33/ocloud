package autonomousdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service provides operations and functionalities related to database management, logging, and compartment handling.
type Service struct {
	repo          domain.AutonomousDatabaseRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance with the provided application context.
func NewService(repo domain.AutonomousDatabaseRepository, appCtx *app.ApplicationContext) *Service {
	return &Service{
		repo:          repo,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}
}

// List retrieves a paginated list of databases with given limit and page number parameters.
// It returns the slice of databases, total count, next page token, and an error if encountered.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]AutonomousDatabase, int, string, error) {
	// Log input parameters at debug level
	logger.LogWithLevel(s.logger, logger.Trace, "FetchPaginatedClusters called with pagination parameters",
		"limit", limit,
		"pageNum", pageNum)

	var databases []AutonomousDatabase
	var nextPageToken string
	var totalCount int

	allDatabases, err := s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to list autonomous databases: %w", err)
	}

	// Apply pagination logic
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= len(allDatabases) {
		logger.LogWithLevel(s.logger, logger.Trace, "Pagination: start index out of bounds", "start", start, "totalDatabases", len(allDatabases))
		return []AutonomousDatabase{}, 0, "", nil // No results for this page
	}

	if end > len(allDatabases) {
		end = len(allDatabases)
		logger.LogWithLevel(s.logger, logger.Trace, "Pagination: adjusted end index", "end", end, "totalDatabases", len(allDatabases))
	}

	databases = make([]AutonomousDatabase, 0, limit)
	for _, db := range allDatabases[start:end] {
		databases = append(databases, AutonomousDatabase(db)) // Convert domain.AutonomousDatabase to local AutonomousDatabase
	}

	totalCount = len(allDatabases)
	if end < len(allDatabases) {
		nextPageToken = "true" // Indicate there's a next page
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

	logger.LogWithLevel(s.logger, logger.Trace, "Completed instance listing with pagination",
		"returnedCount", len(databases),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	logger.Logger.V(logger.Info).Info("Autonomous Database list completed.", "returnedCount", len(databases), "totalCount", totalCount)
	return databases, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search to find autonomous databases matching the given search pattern in their Name field.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]AutonomousDatabase, error) {
	logger.LogWithLevel(s.logger, logger.Trace, "finding database with bleve fuzzy search", "pattern", searchPattern)

	// 1: Fetch all databases
	allDatabases, err := s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all databases: %w", err)
	}
	// 2: Build index
	index, err := util.BuildIndex(allDatabases, func(db domain.AutonomousDatabase) any {
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
			results = append(results, AutonomousDatabase(allDatabases[idx]))
		}
	}
	logger.LogWithLevel(s.logger, logger.Trace, "Compartment search complete", "matches", len(results))

	return results, nil
}

// mapToIndexableDatabase converts an AutonomousDatabase object into an IndexableAutonomousDatabase object.
// It maps only relevant fields required for indexing, such as the database\'s name.
func mapToIndexableDatabase(db domain.AutonomousDatabase) IndexableAutonomousDatabase {
	return IndexableAutonomousDatabase{
		Name: db.Name,
	}
}
