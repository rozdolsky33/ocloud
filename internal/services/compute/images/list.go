package images

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

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
		return fmt.Errorf("listing iamges: %w", err)
	}
	fmt.Println(images, totalCount, nextPageToken)

	return nil
}
