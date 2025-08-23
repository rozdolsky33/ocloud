package image

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociimage "github.com/rozdolsky33/ocloud/internal/oci/compute/image"
)

// ListImages lists all images in the given compartment, allowing the user to select one via a TUI and display its details.
func ListImages(ctx context.Context, appCtx *app.ApplicationContext) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	imageAdapter := ociimage.NewAdapter(computeClient)
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

	var img Image
	for _, it := range images {
		if it.OCID == chosen.Choice() {
			img = it
			break
		}
	}

	err = PrintImageInfo(img, appCtx)
	if err != nil {
		return fmt.Errorf("printing image info: %w", err)
	}

	return nil
}
