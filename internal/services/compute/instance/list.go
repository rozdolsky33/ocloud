package instance

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListInstances lists instances in the configured compartment using the provided application.
// It uses the pre-initialized compute client from the ApplicationContext struct and supports pagination.
func ListInstances(appCtx *app.ApplicationContext, limit int, page int, useJSON bool, showImageDetails bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "ListInstances", "limit", limit, "page", page, "json", useJSON, "showImageDetails", showImageDetails)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	instances, totalCount, nextPageToken, err := service.List(ctx, limit, page, showImageDetails)
	if err != nil {
		return fmt.Errorf("listing instances: %w", err)
	}

	// Display instance information with pagination details
	err = PrintInstancesInfo(instances, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showImageDetails)

	if err != nil {
		return fmt.Errorf("printing instances table: %w", err)
	}

	return nil
}
