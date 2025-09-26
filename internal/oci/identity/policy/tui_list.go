package policy

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewPoliciesListModel constructs a tui.Model titled "Policies" that represents the provided policies,
// mapping each policy to a ResourceItemData with ID, Title, and a Description that combines the policy's
// description and its creation date (formatted as YYYY-MM-DD).
func NewPoliciesListModel(p []domain.Policy) tui.Model {
	return tui.NewModel("Policies", p, func(p domain.Policy) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          p.ID,
			Title:       p.Name,
			Description: fmt.Sprint(p.Description, "  •  ", p.TimeCreated.Format("2006-01-02")),
		}
	})
}
