package loadbalancer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocilb "github.com/rozdolsky33/ocloud/internal/oci/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

func ListLoadBalancers(appCtx *app.ApplicationContext, useJSON, showAll bool) error {
	ctx := context.Background()

	start := time.Now()
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.START", "limit", "page", "json", useJSON, "all")

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

	allLoadBalancers, err := service.ListLoadBalancers(ctx)

	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "list", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("listing load balancers: %w", err)
	}

	//TUI
	model := ocilb.NewLoadBalancerListModel(allLoadBalancers)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting database: %w", err)
	}

	lb, err := service.GetEnrichedLoadBalancer(ctx, id)
	if err != nil {
		return fmt.Errorf("getting load balancer: %w", err)
	}

	return PrintLoadBalancerInfo(lb, appCtx, useJSON, showAll)
}
