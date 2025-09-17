package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/compartment"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetCompartments retrieves and displays a paginated list of compartments.
func GetCompartments(appCtx *app.ApplicationContext, useJSON bool, limit, page int, ocid string) error {
	ctx := context.Background()
	compartmentAdapter := compartment.NewCompartmentAdapter(appCtx.IdentityClient, ocid)
	service := NewService(compartmentAdapter, appCtx.Logger, ocid)

	compartments, totalCount, nextPageToken, err := service.FetchPaginateCompartments(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing compartments: %w", err)
	}

	return PrintCompartmentsTable(compartments, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
}
