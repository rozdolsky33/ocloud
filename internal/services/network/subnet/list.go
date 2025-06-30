package subnet

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func ListSubnets(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing Subnets", "limit", limit, "page", page)

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
	err = PrintSubnetInfo(policies, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)

	if err != nil {
		return fmt.Errorf("printing subnets: %w", err)
	}

	return nil
}
