package networklb

import (
	"context"
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocinlb "github.com/rozdolsky33/ocloud/internal/oci/network/networklb"
)

// SearchNetworkLoadBalancer searches for matching network load balancers and displays their details.
func SearchNetworkLoadBalancer(appCtx *app.ApplicationContext, search string, useJSON, showAll bool) error {
	ctx := context.Background()
	start := time.Now()

	nlbClient, err := oci.NewNetworkLoadBalancerClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.search.error", "stage", "client_init_nlb", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network load balancer client: %w", err)
	}
	nwClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.search.error", "stage", "client_init_network", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network client: %w", err)
	}
	adapter := ocinlb.NewAdapter(nlbClient, nwClient)

	service := NewService(adapter, appCtx)

	matchedNLBs, err := service.FuzzySearch(ctx, search)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.search.error", "stage", "search", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("searching network load balancers: %w", err)
	}

	err = PrintNetworkLoadBalancersInfo(matchedNLBs, appCtx, nil, useJSON, showAll)
	if err != nil {
		return fmt.Errorf("printing network load balancers: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching network load balancers", "search", search, "matched", len(matchedNLBs))
	return nil
}
