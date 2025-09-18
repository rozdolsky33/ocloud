package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	oci_vcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
)

// ListVCNs lists VCNs in the current compartment and prints a table or JSON.
func ListVCNs(appCtx *app.ApplicationContext, useJSON bool) error {
	client, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}
	adapter := oci_vcn.NewAdapter(client)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, err := service.ListVcns(context.Background())
	if err != nil {
		return fmt.Errorf("listing vcns: %w", err)
	}

	return PrintVCNsTable(vcns, appCtx, useJSON)
}
