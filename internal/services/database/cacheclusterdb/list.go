package cacheclusterdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	ocicachecluster "github.com/rozdolsky33/ocloud/internal/oci/database/cacheclusterdb"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListCacheClusters lists all HeatWave Cache Clusters in the application context with TUI.
func ListCacheClusters(appCtx *app.ApplicationContext, useJSON bool) error {
	ctx := context.Background()
	cacheClusterAdapter, err := ocicachecluster.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating cache cluster adapter: %w", err)
	}
	service := NewService(cacheClusterAdapter, appCtx)
	allClusters, err := service.ListCacheClusters(ctx)

	if err != nil {
		return fmt.Errorf("listing cache clusters: %w", err)
	}

	// TUI
	model := ocicachecluster.NewCacheClusterListModel(allClusters)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting cache cluster: %w", err)
	}

	cluster, err := service.repo.GetCacheCluster(ctx, id)
	if err != nil {
		return fmt.Errorf("getting cache cluster: %w", err)
	}

	return PrintCacheClusterInfo(cluster, appCtx, useJSON, true)
}
