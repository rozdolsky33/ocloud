package image

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewImageListModel builds a TUI list for images.
func NewImageListModel(images []domain.Image) tui.Model {
	return tui.NewModel("Images", images, func(image domain.Image) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          image.OCID,
			Title:       image.DisplayName,
			Description: fmt.Sprintf("OS: %s • Version: %s", image.OperatingSystem, image.OperatingSystemVersion),
		}
	})
}
