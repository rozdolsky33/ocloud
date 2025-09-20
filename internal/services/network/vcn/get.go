package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociVcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetVCNs retrieves a VCN by OCID and prints its summary or JSON.
func GetVCNs(appCtx *app.ApplicationContext, limit, page int, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	adapter := ociVcn.NewAdapter(networkClient)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	vcns, totalCount, nextPageToken, err := service.FetchPaginatedVCNs(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("getting vcn: %w", err)
	}

	return PrintVCNsInfo(vcns, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
