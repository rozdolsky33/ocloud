package image

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociimage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
)

// FindImages finds and displays images matching a name pattern.
func FindImages(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	imageAdapter := ociimage.NewAdapter(computeClient)
	service := NewService(imageAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedImages, err := service.Find(context.Background(), namePattern)
	if err != nil {
		return fmt.Errorf("finding images: %w", err)
	}

	return PrintImagesInfo(matchedImages, appCtx, nil, useJSON)
}
