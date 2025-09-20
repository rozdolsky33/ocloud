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

// ListVCNs lists VCNs in the current compartment and prints a table or JSON.
func ListVCNs(appCtx *app.ApplicationContext, useJSON bool) error {
	client, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}
	adapter := ociVcn.NewAdapter(client)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, err := service.ListVcns(context.Background())
	if err != nil {
		return fmt.Errorf("listing vcns: %w", err)
	}

	//TUI
	model := ociVcn.NewVCNListModel(vcns)
	id, err := listx.Run(model)
	if err != nil {
		if errors.Is(err, listx.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting vcn: %w", err)
	}

	vcn, err := service.vcnRepo.GetVcn(context.Background(), id)
	if err != nil {
		return fmt.Errorf("getting vcn: %w", err)
	}

	return PrintVCNSummary(vcn, appCtx, useJSON)
}
