package image

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
)

func ListImages(ctx context.Context, appCtx *app.ApplicationContext) error {

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating image service: %w", err)
	}
	images, err := service.fetchAllImages(ctx)

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
		if it.ID == chosen.Choice() {
			img = it
			break
		}
	}

	fmt.Println(img)

	return nil
}
