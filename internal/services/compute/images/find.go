package images

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func FindImages(appCtx *app.ApplicationContext, searchPattern string, useJSON bool) error {

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(appCtx.Logger, 1, "FindImage", "json", useJSON)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}
	ctx := context.Background()
	matchedImages, err := service.Find(ctx, searchPattern)
	if err != nil {
		return fmt.Errorf("finding images: %w", err)
	}
	err = PrintImagesInfo(matchedImages, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing images table: %w", err)
	}

	return nil
}
