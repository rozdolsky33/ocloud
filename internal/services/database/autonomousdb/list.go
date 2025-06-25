package autonomousdb

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListAutonomousDatabase fetches and lists all autonomous databases within a specified application context.
// appCtx represents the application context containing configuration and client details.
// useJSON if true, outputs the list of databases in JSON format; otherwise, uses a plain-text format.
// Returns an error if the operation fails.
func ListAutonomousDatabase(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing Autonomous Databases")

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating autonomous database service: %w", err)
	}

	ctx := context.Background()
	allDatabases, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing autonomous databases: %w", err)
	}

	// Display image information with pagination details
	err = PrintAutonomousDbInfo(allDatabases, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)

	if err != nil {
		return fmt.Errorf("printing image table: %w", err)
	}

	return nil
}
