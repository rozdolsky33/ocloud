package policy

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

func NewPoliciesListModel(p []domain.Policy) listx.Model {
	return listx.NewModel("Policies", p, func(p domain.Policy) listx.ResourceItemData {
		return listx.ResourceItemData{
			ID:          p.ID,
			Title:       p.Name,
			Description: fmt.Sprintf("Description: %s  •  Created: %s", p.Description, p.TimeCreated.Format("2006-01-02")),
		}
	})
}
