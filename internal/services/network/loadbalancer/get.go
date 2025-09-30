package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	oci "github.com/rozdolsky33/ocloud/internal/oci"
	ocilb "github.com/rozdolsky33/ocloud/internal/oci/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetLoadBalancers retrieves load balancers and displays a paginated list.
func GetLoadBalancers(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	start := time.Now()
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.START", "limit", limit, "page", page, "json", useJSON, "all", showAll)

	lbClient, err := oci.NewLoadBalancerClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "client_init_lb", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating load balancer client: %w", err)
	}
	nwClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "client_init_network", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network client: %w", err)
	}
	certsClient, err := oci.NewCertificatesManagementClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "client_init_certs", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating certificates management client: %w", err)
	}
	adapter := ocilb.NewAdapter(lbClient, nwClient, certsClient)

	service := NewService(adapter, appCtx)

	ctx := context.Background()
	lbs, totalCount, nextPageToken, err := service.FetchPaginatedLoadBalancers(ctx, limit, page, showAll)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "fetch", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("listing load balancers: %w", err)
	}

	logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.FINISH", "count", len(lbs), "total_count", totalCount, "next_page", nextPageToken, "duration_ms", time.Since(start).Seconds())
	return PrintLoadBalancersInfo(lbs, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showAll)
}
