package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociadb "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetAutonomousDatabase retrieves a list of Autonomous Databases and displays them in a table or JSON format.
func GetAutonomousDatabase(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing Autonomous Databases")
	adapter, err := ociadb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating database adapter: %w", err)
	}

	service := NewService(adapter, appCtx)

	ctx := context.Background()
	allDatabases, totalCount, nextPageToken, err := service.FetchPaginatedAutonomousDb(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing autonomous databases: %w", err)
	}

	return PrintAutonomousDbInfo(allDatabases, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll)
}
