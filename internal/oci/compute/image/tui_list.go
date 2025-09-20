package image

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

// NewImageListModel builds a TUI list for images.
func NewImageListModel(images []domain.Image) listx.Model {
	return listx.NewModel("Images", images, func(image domain.Image) listx.ResourceItemData {
		return listx.ResourceItemData{
			ID:          image.OCID,
			Title:       image.DisplayName,
			Description: fmt.Sprintf("OS: %s  •  Version: %s", image.OperatingSystem, image.OperatingSystemVersion),
		}
	})
}
