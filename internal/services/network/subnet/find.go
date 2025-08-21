package subnet

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocisubnet "github.com/rozdolsky33/ocloud/internal/oci/network/subnet"
)

// FindSubnets finds and displays subnets matching a name pattern.
func FindSubnets(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	subnetAdapter := ocisubnet.NewAdapter(networkClient)
	service := NewService(subnetAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedSubnets, err := service.Find(context.Background(), namePattern)
	if err != nil {
		return fmt.Errorf("finding subnets: %w", err)
	}

	return PrintSubnetInfo(matchedSubnets, appCtx, useJSON)
}
