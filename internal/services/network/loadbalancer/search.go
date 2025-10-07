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

// SearchLoadBalancer searches for matching load balancers based on a fuzzy search string and displays their details.
// appCtx provides context and clients for API calls.
// search specifies the fuzzy search string to filter load balancers.
// useJSON determines if the output should be in JSON format.
// showAll includes all details about load balancers in the output if set to true.
// Returns an error if there is a failure in the process.
func SearchLoadBalancer(appCtx *app.ApplicationContext, search string, useJSON, showAll bool) error {
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

	matchedLoadBalancers, err := service.FuzzySearch(ctx, search)

	if err != nil {
		logger.LogWithLevel(appCtx.Logger, logger.Debug, "lb.service.get.error", "stage", "list", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("listing load balancers: %w", err)
	}

	err = PrintLoadBalancersInfo(matchedLoadBalancers, appCtx, nil, useJSON, showAll)
	if err != nil {
		return fmt.Errorf("printing load balancers: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching load balancers", "search", search, "matched", len(matchedLoadBalancers))
	return nil
}
