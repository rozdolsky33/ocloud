package compartment

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewPoliciesListModel builds a TUI list model titled "Compartments" from the provided slice of domain.Compartment.
// Each compartment is mapped to a ResourceItemData with ID set to OCID, Title set to DisplayName, and Description
// composed of the LifecycleState, a bullet, and the compartment Description.
func NewPoliciesListModel(c []domain.Compartment) tui.Model {
	return tui.NewModel("Compartments", c, func(c domain.Compartment) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          c.OCID,
			Title:       c.DisplayName,
			Description: fmt.Sprint(c.LifecycleState, "  •  ", c.Description),
		}
	})
}
