package subnet

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocisubnet "github.com/rozdolsky33/ocloud/internal/oci/network/subnet"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListSubnets retrieves and displays a paginated list of subnets.
func ListSubnets(appCtx *app.ApplicationContext, useJSON bool, limit, page int, sortBy string) error {
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	subnetAdapter := ocisubnet.NewAdapter(networkClient)
	service := NewService(subnetAdapter, appCtx.Logger, appCtx.CompartmentID)

	subnets, totalCount, nextPageToken, err := service.List(context.Background(), limit, page)
	if err != nil {
		return fmt.Errorf("listing subnets: %w", err)
	}

	return PrintSubnetTable(subnets, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, sortBy)
}
