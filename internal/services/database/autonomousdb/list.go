package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ocidbadapter "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListAutonomousDatabase fetches and lists all autonomous databases within a specified application context.
// appCtx represents the application context containing configuration and client details.
// useJSON if true, outputs the list of databases in JSON format; otherwise, uses a plain-text format.
func ListAutonomousDatabase(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing Autonomous Databases")

	adapter, err := ocidbadapter.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating database adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	allDatabases, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing autonomous databases: %w", err)
	}

	// Convert to a domain type for printing
	domainDbs := make([]domain.AutonomousDatabase, 0, len(allDatabases))
	for _, db := range allDatabases {
		domainDbs = append(domainDbs, domain.AutonomousDatabase(db))
	}

	// Display database information with pagination details
	if err := PrintAutonomousDbInfo(domainDbs, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON); err != nil {
		return fmt.Errorf("printing image table: %w", err)
	}
	return nil
}
