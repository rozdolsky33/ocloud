package subnet

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListSubnets retrieves and displays a list of subnets, utilizing pagination, sorting, and optional JSON output.
// Parameters include an application context, a flag for JSON output, limit, page number, and a sorting key.
// Returns an error if any issue occurs during retrieval or output generation.
func ListSubnets(appCtx *app.ApplicationContext, useJSON bool, limit, page int, sortBy string) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing Subnets", "limit", limit, "page", page, "sortBy", sortBy)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating subnet service: %w", err)
	}

	ctx := context.Background()
	policies, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing subnets: %w", err)
	}

	// Display policies information with pagination details
	err = PrintSubnetTable(policies, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, sortBy)

	if err != nil {
		return fmt.Errorf("printing subnets: %w", err)
	}

	return nil
}
