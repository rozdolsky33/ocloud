package dynamicgroup

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewDynamicGroupListModel builds a TUI list for dynamic groups.
func NewDynamicGroupListModel(dgs []domain.DynamicGroup) tui.Model {
	return tui.NewModel("Dynamic Groups", dgs, func(dg domain.DynamicGroup) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          dg.OCID,
			Title:       dg.Name,
			Description: fmt.Sprint(dg.LifecycleState, "  •  ", dg.Description),
		}
	})
}
