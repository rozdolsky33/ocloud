package vcn

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociVcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
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

	model := ociVcn.NewVCNListModel(vcns)
	id, err := listx.Run(model)
	if err != nil {
		if errors.Is(err, listx.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("listing vcn: %w", err)
	}

	vcn, err := service.vcnRepo.GetEnrichedVcn(ctx, id)
	if err != nil {
		return fmt.Errorf("getting vcn: %w", err)
	}

	return PrintVCNInfo(vcn, appCtx, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
