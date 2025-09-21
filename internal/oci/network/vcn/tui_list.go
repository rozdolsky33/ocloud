package vcn

import (
	"fmt"

	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewVCNListModel builds a TUI list for VCNs.
func NewVCNListModel(v []domain.VCN) tui.Model {
	return tui.NewModel("VCNs", v, func(v domain.VCN) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          v.OCID,
			Title:       v.DisplayName,
			Description: fmt.Sprint(v.LifecycleState, " • ", v.DomainName, " • ", v.TimeCreated.Format("2006-01-02")),
		}
	})
}
