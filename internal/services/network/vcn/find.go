package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocivcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
)

func FindVCNs(appCtx *app.ApplicationContext, pattern string, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	adapter := ocivcn.NewAdapter(networkClient)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, err := service.Find(ctx, pattern)
	if err != nil {
		return fmt.Errorf("finding vcn: %w", err)
	}

	return PrintVCNsInfo(vcns, appCtx, nil, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
