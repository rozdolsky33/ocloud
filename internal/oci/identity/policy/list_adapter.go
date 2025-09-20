package policy

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

// NewPoliciesListModel builds a TUI list for policies.
func NewPoliciesListModel(p []domain.Policy) listx.Model {
	return listx.NewModel("Policies", p, func(p domain.Policy) listx.ResourceItemData {
		return listx.ResourceItemData{
			ID:          p.ID,
			Title:       p.Name,
			Description: fmt.Sprint(p.Description, "  •  ", p.TimeCreated.Format("2006-01-02")),
		}
	})
}
