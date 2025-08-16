package image

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindImages retrieves image matching the provided search pattern and outputs them in either table or JSON format.
// Parameters:
// - appCtx: The application context containing configuration, clients, and logging.
// - searchPattern: The string used to search for matching image.
// - useJSON: A boolean indicating whether to output the result in JSON format.
// Returns an error if the operation fails at any stage (service creation, image retrieval, or output).
func FindImages(appCtx *app.ApplicationContext, searchPattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "FindImage", "json", useJSON)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	matchedImages, err := service.Find(ctx, searchPattern)
	if err != nil {
		return fmt.Errorf("finding image: %w", err)
	}

	err = PrintImagesInfo(matchedImages, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing image table: %w", err)
	}

	return nil
}
