package image

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociImage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
)

// ListImages lists all images in the given compartment, allowing the user to select one via a TUI and display its details.
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

	// TUI selection
	im := NewImageListModelFancy(images)
	ip := tea.NewProgram(im, tea.WithContext(ctx))
	ires, err := ip.Run()
	if err != nil {
		return fmt.Errorf("image selection TUI: %w", err)
	}
	chosen, ok := ires.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return err
	}

	image, err := service.imageRepo.GetImage(ctx, chosen.Choice())
	if err != nil {
		return fmt.Errorf("getting image: %w", err)
	}

	return PrintImageInfo(image, appCtx, useJSON)
}
