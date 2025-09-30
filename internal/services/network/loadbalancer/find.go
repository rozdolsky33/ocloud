package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocilb "github.com/rozdolsky33/ocloud/internal/oci/network/loadbalancer"
)

func FindLoadBalancer(appCtx *app.ApplicationContext, namePattern string, useJSON, showAll bool) error {
	ctx := context.Background()
	start := time.Now()

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

	matchedLoadBalancers, err := service.Find(ctx, namePattern)

	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "list", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("listing load balancers: %w", err)
	}

	return PrintLoadBalancersInfo(matchedLoadBalancers, appCtx, nil, useJSON, showAll)
}
