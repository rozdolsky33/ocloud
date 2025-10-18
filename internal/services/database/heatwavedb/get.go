package heatwavedb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociheatwave "github.com/rozdolsky33/ocloud/internal/oci/database/heatwavedb"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetHeatWaveDatabase retrieves a list of HeatWave Databases and displays them in a table or JSON format.
func GetHeatWaveDatabase(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing HeatWave Databases")
	adapter, err := ociheatwave.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating HeatWave database adapter: %w", err)
	}

	service := NewService(adapter, appCtx)

	ctx := context.Background()
	allDatabases, totalCount, nextPageToken, err := service.FetchPaginatedHeatWaveDb(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing HeatWave databases: %w", err)
	}

	return PrintHeatWaveDbsInfo(allDatabases, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll)
}
