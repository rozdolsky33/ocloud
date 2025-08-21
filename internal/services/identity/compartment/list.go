package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListCompartments retrieves and displays a paginated list of compartments.
func ListCompartments(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	appCtx.Logger.V(1).Info("listing compartments", "limit", limit, "page", page)

	// Create the infrastructure adapter.
	compartmentAdapter := identity.NewCompartmentAdapter(appCtx.IdentityClient, appCtx.TenancyID)

	// Create the application service, injecting the adapter.
	// The service is now decoupled from the OCI SDK.
	service := NewService(compartmentAdapter, appCtx.Logger, appCtx.TenancyID)

	ctx := context.Background()
	compartments, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing compartments: %w", err)
	}

	// Display compartment information with pagination details.
	err = PrintCompartmentsTable(compartments, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)

	if err != nil {
		return fmt.Errorf("printing compartments: %w", err)
	}

	return nil
}
