package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocivcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
)

// FindVCNs finds virtual cloud networks in the application's compartment that match the given pattern and prints their information.
// The `pattern` filters VCNs by name or OCID. The boolean flags control output: `useJSON` emits JSON, and `gateways`, `subnets`, `nsgs`, `routes`, `securityLists` include details for each respective resource in the output.
// It returns an error encountered while creating the network client, searching for VCNs, or printing the results.
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
