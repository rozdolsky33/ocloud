package vcn

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocivcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListVCNs lists VCNs in the configured compartment, opens a TUI to select one, fetches enriched details for the selected VCN, and prints its information according to the provided flags.
// If the TUI is cancelled by the user, ListVCNs returns nil. It returns an error when creating the network client, retrieving the VCN list, fetching the enriched VCN, or printing the information.
func ListVCNs(appCtx *app.ApplicationContext, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	adapter := ocivcn.NewAdapter(networkClient)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, err := service.ListVcns(ctx)
	if err != nil {
		return fmt.Errorf("getting vcn: %w", err)
	}

	model := ocivcn.NewVCNListModel(vcns)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
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
