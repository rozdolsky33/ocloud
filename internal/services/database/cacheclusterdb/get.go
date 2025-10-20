package cacheclusterdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ocicachecluster "github.com/rozdolsky33/ocloud/internal/oci/database/cacheclusterdb"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetCacheClusters retrieves a list of HeatWave Cache Clusters and displays them in a table or JSON format.
func GetCacheClusters(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing HeatWave Cache Clusters")
	adapter, err := ocicachecluster.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating cache cluster adapter: %w", err)
	}

	service := NewService(adapter, appCtx)

	ctx := context.Background()
	allClusters, totalCount, nextPageToken, err := service.FetchPaginatedCacheClusters(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing cache clusters: %w", err)
	}

	return PrintCacheClustersInfo(allClusters, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll)
}
