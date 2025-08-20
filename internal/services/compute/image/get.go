package image

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetImages retrieves and displays a list of images, utilizing pagination, sorting, and optional JSON output.
// It uses the application context for accessing configurations, logging, and output streams.
func GetImages(appCtx *app.ApplicationContext, limit int, page int, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "ListImages", "limit", limit, "page", page, "json", useJSON)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	images, totalCount, nextPageToken, err := service.Get(ctx, limit, page)
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
