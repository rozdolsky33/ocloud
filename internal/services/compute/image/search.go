package image

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociImage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
)

// SearchImages performs a search for images based on a given search term.
// It uses fuzzy matching to find relevant images in the specified OCI compartment.
// The results are printed in either a tabular or JSON format depending on the `useJSON` flag.
// An error is returned if there are issues with creating required clients, searching images, or printing results.
func SearchImages(appCtx *app.ApplicationContext, search string, useJSON bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	imageAdapter := ociImage.NewAdapter(computeClient)
	service := NewService(imageAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedImages, err := service.FuzzySearch(context.Background(), search)
	if err != nil {
		return fmt.Errorf("finding images: %w", err)
	}

	err = PrintImagesInfo(matchedImages, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing images: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching images", "search", search, "matched", len(matchedImages))
	return nil
}
