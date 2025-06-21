package image

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListImages lists image from the compute service with provided limit, page, and JSON output option.
// It uses the application context for configuration and logging.
// Returns an error if the operation fails.
func ListImages(appCtx *app.ApplicationContext, limit int, page int, useJSON bool) error {
	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(appCtx.Logger, 1, "ListImages", "limit", limit, "page", page, "json", useJSON)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	images, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing images: %w", err)
	}

	// Display image information with pagination details
	err = PrintImagesInfo(images, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
	if err != nil {
		return fmt.Errorf("printing image table: %w", err)
	}

	return nil
}
