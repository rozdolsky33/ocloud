package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListCompartments retrieves and displays a paginated list of compartments for the given application context.
// It uses the provided limit, page number, and output format (JSON or table) to render the results.
func ListCompartments(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing Compartments", "limit", limit, "page", page)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compartment service: %w", err)
	}

	ctx := context.Background()
	compartments, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing compartments: %w", err)
	}

	// Display compartment information with pagination details
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
