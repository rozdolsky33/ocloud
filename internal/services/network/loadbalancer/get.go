package loadbalancer

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	oci "github.com/rozdolsky33/ocloud/internal/oci"
	ocilb "github.com/rozdolsky33/ocloud/internal/oci/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetLoadBalancers retrieves load balancers and displays a paginated list.
func GetLoadBalancers(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing Load Balancers")
	lbClient, err := oci.NewLoadBalancerClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating load balancer client: %w", err)
	}
	nwClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}
	certsClient, err := oci.NewCertificatesManagementClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating certificates management client: %w", err)
	}
	adapter := ocilb.NewAdapter(lbClient, nwClient, certsClient)

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
