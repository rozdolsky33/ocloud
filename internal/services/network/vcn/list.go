package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociVcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
)

func ListVCNs(appCtx *app.ApplicationContext, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	adapter := ociVcn.NewAdapter(networkClient)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, err := service.ListVcns(ctx)
	if err != nil {
		return fmt.Errorf("getting vcn: %w", err)
	}

	return PrintVCNsInfo(vcns, appCtx, nil, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
