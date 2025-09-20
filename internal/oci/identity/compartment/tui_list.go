package compartment

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

// NewPoliciesListModel builds a TUI list for policies.
func NewPoliciesListModel(c []domain.Compartment) listx.Model {
	return listx.NewModel("Compartments", c, func(c domain.Compartment) listx.ResourceItemData {
		return listx.ResourceItemData{
			ID:          c.OCID,
			Title:       c.DisplayName,
			Description: fmt.Sprint(c.LifecycleState, "  •  ", c.Description),
		}
	})
}
