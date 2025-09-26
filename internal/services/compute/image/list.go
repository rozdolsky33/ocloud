package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociImage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListImages lists images in the application's compartment, presents a TUI for the user to select an image, and prints the selected image's details.
// It returns nil if the user cancels selection or after successfully printing the selected image; it returns a wrapped error when creating the compute client,
// listing images, selecting an image, or retrieving the chosen image fails.
func ListImages(ctx context.Context, appCtx *app.ApplicationContext, useJSON bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	imageAdapter := ociImage.NewAdapter(computeClient)
	service := NewService(imageAdapter, appCtx.Logger, appCtx.CompartmentID)

	images, err := service.imageRepo.ListImages(ctx, appCtx.CompartmentID)
	if err != nil {
		return fmt.Errorf("listing images: %w", err)
	}

	// TUI
	model := ociImage.NewImageListModel(images)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting image: %w", err)
	}

	image, err := service.imageRepo.GetImage(ctx, id)
	if err != nil {
		return fmt.Errorf("getting image: %w", err)
	}

	return PrintImageInfo(image, appCtx, useJSON)
}
