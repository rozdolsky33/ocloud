package networklb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocinlb "github.com/rozdolsky33/ocloud/internal/oci/network/networklb"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

func ListNetworkLoadBalancers(appCtx *app.ApplicationContext, useJSON, showAll bool) error {
	ctx := context.Background()

	start := time.Now()
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.list.START", "json", useJSON, "all", showAll)

	nlbClient, err := oci.NewNetworkLoadBalancerClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.list.error", "stage", "client_init_nlb", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network load balancer client: %w", err)
	}
	nwClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.list.error", "stage", "client_init_network", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("creating network client: %w", err)
	}
	adapter := ocinlb.NewAdapter(nlbClient, nwClient)

	service := NewService(adapter, appCtx)

	allNLBs, err := service.ListNetworkLoadBalancers(ctx)
	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "nlb.service.list.error", "stage", "list", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("listing network load balancers: %w", err)
	}

	// TUI
	model := ocinlb.NewNetworkLoadBalancerListModel(allNLBs)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting network load balancer: %w", err)
	}

	nlb, err := service.GetEnrichedNetworkLoadBalancer(ctx, id)
	if err != nil {
		return fmt.Errorf("getting network load balancer: %w", err)
	}

	return PrintNetworkLoadBalancerInfo(nlb, appCtx, useJSON, showAll)
}
