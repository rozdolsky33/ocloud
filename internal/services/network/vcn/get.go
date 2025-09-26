package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocivcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetVCNs retrieves VCNs for the application's compartment, fetches a paginated
// list, and prints their information according to the requested output flags.
// It initializes an OCI network client from the application context, obtains
// paginated VCN results, and delegates formatting/printing (summary or JSON and
// optional details: gateways, subnets, NSGs, routes, security lists).
// It returns an error if client initialization, VCN retrieval, or printing fails.
func GetVCNs(appCtx *app.ApplicationContext, limit, page int, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	adapter := ocivcn.NewAdapter(networkClient)
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
