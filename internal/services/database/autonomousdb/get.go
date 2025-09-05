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

// GetAutonomousDatabase retrieves a list of Autonomous Databases and displays them in a table or JSON format.
func GetAutonomousDatabase(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
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

	// Convert to a domain type and best-effort enriches each item with a full Get call
	domainDbs := make([]domain.AutonomousDatabase, 0, len(allDatabases))
	for _, db := range allDatabases {
		basic := domain.AutonomousDatabase(db)
		full, gerr := adapter.GetAutonomousDatabase(ctx, basic.ID)
		if gerr == nil && full != nil {
			domainDbs = append(domainDbs, *full)
		} else {
			domainDbs = append(domainDbs, basic)
		}
	}

	// Display database information with pagination details
	if err := PrintAutonomousDbInfo(domainDbs, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll); err != nil {
		return fmt.Errorf("printing image table: %w", err)
	}
	return nil
}
