package image

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociimage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetImages retrieves and displays a paginated list of images.
func GetImages(appCtx *app.ApplicationContext, limit int, page int, useJSON bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	imageAdapter := ociimage.NewAdapter(computeClient)
	service := NewService(imageAdapter, appCtx.Logger, appCtx.CompartmentID)

	images, totalCount, nextPageToken, err := service.Get(context.Background(), limit, page)
	if err != nil {
		return fmt.Errorf("listing images: %w", err)
	}

	return PrintImagesInfo(images, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
}
