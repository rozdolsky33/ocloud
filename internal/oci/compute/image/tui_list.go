package image

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewImageListModel creates a TUI list model titled "Images" from the provided images.
// Each image is mapped to a ResourceItemData with ID set to the image OCID, Title set to the display name,
// and Description formatted as "OS: <OperatingSystem>  •  Version: <OperatingSystemVersion>".
func NewImageListModel(images []domain.Image) tui.Model {
	return tui.NewModel("Images", images, func(image domain.Image) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          image.OCID,
			Title:       image.DisplayName,
			Description: fmt.Sprintf("OS: %s  •  Version: %s", image.OperatingSystem, image.OperatingSystemVersion),
		}
	})
}
