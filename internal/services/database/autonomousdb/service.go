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

	allDatabases, err := s.repo.ListAutonomousDatabases(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to list autonomous databases: %w", err)
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
			results = append(results, allDatabases[idx]) // TODO: REDO
		}
	}
	return results, nil
}

func mapToIndexableDatabase(db domain.AutonomousDatabase) IndexableAutonomousDatabase {
	return IndexableAutonomousDatabase{
		Name: db.Name,
	}
}
