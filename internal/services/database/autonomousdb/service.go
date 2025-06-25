package autonomousdb

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"strings"
)

// NewService initializes a new Service instance with the provided application context.
// Returns a Service pointer and an error if initialization fails.
func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	cfg := appCtx.Provider
	dbClient, err := oci.NewDatabaseClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("faild to create database client: %w", err)
	}
	return &Service{
		dbClient:      dbClient,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}, nil
}

// List retrieves a paginated list of databases with given limit and page number parameters.
// It returns the slice of databases, total count, next page token, and an error if encountered.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]AutonomousDatabase, int, string, error) {
	// Log input parameters at debug level
	logger.LogWithLevel(s.logger, 3, "List called with pagination parameters",
		"limit", limit,
		"pageNum", pageNum)

	var allDatabases []AutonomousDatabase
	var nextPageToken string
	var totalCount int

	// Create a request with a limit parameter to fetch only the required page
	request := database.ListAutonomousDatabasesRequest{
		CompartmentId: &s.compartmentID,
	}

	if limit > 0 {
		request.Limit = &limit
		logger.LogWithLevel(s.logger, 3, "Setting limit parameter", "limit", limit)
	}

	// Add limit parameters if specified
	if pageNum > 1 && limit > 0 {
		logger.LogWithLevel(s.logger, 3, "Calculating page token for page", "pageNum", pageNum)

		// We need to fetch page tokens until we reach the desired page
		page := ""
		currentPage := 1

		for currentPage < pageNum {
			// Fetch just the page token, not actual data
			// Use the same limit to ensure consistent pagination
			tokenRequest := database.ListAutonomousDatabasesRequest{
				CompartmentId: &s.compartmentID,
				Page:          &page,
			}
			if limit > 0 {
				tokenRequest.Limit = &limit
			}

			resp, err := s.dbClient.ListAutonomousDatabases(ctx, tokenRequest)
			if err != nil {
				return nil, 0, "", fmt.Errorf("fetching page token: %w", err)
			}

			// If there's no next page, we've reached the end
			if resp.OpcNextPage == nil {
				logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return an empty result since the requested page is beyond available data
				return []AutonomousDatabase{}, 0, "", nil
			}
			// Move to the next page
			page = *resp.OpcNextPage
			currentPage++
		}
		// Set the page toke for the actual request
		request.Page = &page
		logger.LogWithLevel(s.logger, 3, "Using page token for page", "pageNum", pageNum, "token", page)
	}

	// Fetch database for the request
	resp, err := s.dbClient.ListAutonomousDatabases(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing database: %w", err)
	}
	// Set the total count to the number of instances returned
	// If we have a next page, this is an estimate
	totalCount = len(resp.Items)
	//If we have a next page, we know there are more instances
	if resp.OpcNextPage != nil {
		// Estimate total count based on current page and items per rage
		totalCount = pageNum*limit + limit
	}

	// Process the databases
	for _, item := range resp.Items {
		db := mapToDatabase(item)
		allDatabases = append(allDatabases, db)
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

	logger.LogWithLevel(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(allDatabases),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	return allDatabases, totalCount, nextPageToken, nil
}

func (s *Service) Find(ctx context.Context, searchPattern string) ([]AutonomousDatabase, error) {
	logger.LogWithLevel(s.logger, 3, "finding database with bleve fuzzy search", "pattern", searchPattern)

	// 1: Fetch all databases
	allDatabases, err := s.fetchAllAutonomousDatabases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all databases: %w", err)
	}
	// 2: Build index
	index, err := util.BuildIndex(allDatabases, func(db AutonomousDatabase) any {
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
			results = append(results, allDatabases[idx])
		}
	}
	logger.LogWithLevel(s.logger, 2, "Compartment search complete", "matches", len(results))

	return results, nil
}

func (s *Service) fetchAllAutonomousDatabases(ctx context.Context) ([]AutonomousDatabase, error) {
	var allDatabases []AutonomousDatabase
	page := ""
	for {
		resp, err := s.dbClient.ListAutonomousDatabases(ctx, database.ListAutonomousDatabasesRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list database: %w", err)
		}
		for _, item := range resp.Items {
			db := mapToDatabase(item)
			allDatabases = append(allDatabases, db)
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return allDatabases, nil
}

func mapToIndexableDatabase(db AutonomousDatabase) IndexableAutonomousDatabase {
	return IndexableAutonomousDatabase{
		Name: db.Name,
	}
}

func mapToDatabase(db database.AutonomousDatabaseSummary) AutonomousDatabase {
	return AutonomousDatabase{
		Name:              *db.DbName,
		ID:                *db.Id,
		PrivateEndpoint:   *db.PrivateEndpoint,
		PrivateEndpointIp: *db.PrivateEndpointIp,
		ConnectionStrings: db.ConnectionStrings.AllConnectionStrings,
		Profiles:          db.ConnectionStrings.Profiles,
	}
}
