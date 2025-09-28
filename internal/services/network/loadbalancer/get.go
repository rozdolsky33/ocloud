package loadbalancer

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ocilb "github.com/rozdolsky33/ocloud/internal/oci/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetLoadBalancers retrieves load balancers and displays a paginated list.
func GetLoadBalancers(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing Load Balancers")
	adapter, err := ocilb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating load balancer adapter: %w", err)
	}

	service := NewService(adapter, appCtx)

	ctx := context.Background()
	lbs, totalCount, nextPageToken, err := service.FetchPaginatedLoadBalancers(ctx, limit, page, showAll)
	if err != nil {
		return fmt.Errorf("listing load balancers: %w", err)
	}

	return PrintLoadBalancersInfo(lbs, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll)
}
