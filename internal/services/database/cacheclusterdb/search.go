package cacheclusterdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ocicachecluster "github.com/rozdolsky33/ocloud/internal/oci/database/cacheclusterdb"
)

// SearchCacheClusters searches for OCI HeatWave Cache Clusters matching the given query string in the current context.
func SearchCacheClusters(appCtx *app.ApplicationContext, search string, useJSON bool, showAll bool) error {
	adapter, err := ocicachecluster.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating cache cluster adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	matchedClusters, err := service.FuzzySearch(ctx, search)
	if err != nil {
		return fmt.Errorf("finding cache clusters: %w", err)
	}
	err = PrintCacheClustersInfo(matchedClusters, appCtx, nil, useJSON, showAll)
	if err != nil {
		return fmt.Errorf("printing cache clusters: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching cache clusters", "search", search, "matched", len(matchedClusters))
	return nil
}
