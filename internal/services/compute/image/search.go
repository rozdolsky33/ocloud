package image

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociImage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
)

// SearchImages finds and displays images matching a pattern using fuzzy/prefix/substring search over indexed fields.
func SearchImages(appCtx *app.ApplicationContext, searchPattern string, useJSON bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	imageAdapter := ociImage.NewAdapter(computeClient)
	service := NewService(imageAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedImages, err := service.FuzzySearch(context.Background(), searchPattern)
	if err != nil {
		return fmt.Errorf("finding images: %w", err)
	}

	err = PrintImagesInfo(matchedImages, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing images: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching images", "searchPattern", searchPattern, "matched", len(matchedImages))
	return nil
}
