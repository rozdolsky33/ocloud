package networklb

import (
	"context"
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	oci "github.com/rozdolsky33/ocloud/internal/oci"
	ocinlb "github.com/rozdolsky33/ocloud/internal/oci/network/networklb"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetNetworkLoadBalancers retrieves network load balancers and displays a paginated list.
func GetNetworkLoadBalancers(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	start := time.Now()
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.get.START", "limit", limit, "page", page, "json", useJSON, "all", showAll)

	nlbClient, err := oci.NewNetworkLoadBalancerClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.get.error", "stage", "client_init_nlb", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network load balancer client: %w", err)
	}
	nwClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.get.error", "stage", "client_init_network", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network client: %w", err)
	}
	adapter := ocinlb.NewAdapter(nlbClient, nwClient)

	service := NewService(adapter, appCtx)

	ctx := context.Background()
	nlbs, totalCount, nextPageToken, err := service.FetchPaginatedNetworkLoadBalancers(ctx, limit, page, showAll)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.get.error", "stage", "fetch", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("listing network load balancers: %w", err)
	}

	logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.get.FINISH", "count", len(nlbs), "total_count", totalCount, "next_page", nextPageToken, "duration_ms", time.Since(start).Seconds())
	return PrintNetworkLoadBalancersInfo(nlbs, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll)
}
