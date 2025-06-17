package instance

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// ListInstances lists instances in the configured compartment using the provided application.
// It uses the pre-initialized compute client from the AppContext struct and supports pagination.
func ListInstances(appCtx *app.AppContext, limit int, page int, useJSON bool) error {
	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(appCtx.Logger, 1, "ListInstances()", "limit", limit, "page", page, "json", useJSON)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	instances, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing instances: %w", err)
	}

	// Display instance information with pagination details
	PrintInstancesTable(instances, appCtx, &PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)

	return nil
}
