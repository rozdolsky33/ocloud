package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocivcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
)

func SearchVCNs(appCtx *app.ApplicationContext, pattern string, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	adapter := ocivcn.NewAdapter(networkClient)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, err := service.FuzzySearch(ctx, pattern)
	if err != nil {
		return fmt.Errorf("finding vcn: %w", err)
	}
	err = PrintVCNsInfo(vcns, appCtx, nil, useJSON, gateways, subnets, nsgs, routes, securityLists)
	if err != nil {
		return fmt.Errorf("printing vcn: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching vcn", "search", pattern, "matched", len(vcns))
	return nil
}
