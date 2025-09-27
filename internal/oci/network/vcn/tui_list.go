package vcn

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewVCNListModel builds a TUI list for VCNs.
func NewVCNListModel(v []domain.VCN) tui.Model {
	return tui.NewModel("VCNs", v, func(v domain.VCN) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          v.OCID,
			Title:       v.DisplayName,
			Description: describeVCN(v),
		}
	})
}

// describeVCN constructs a concise description of a VCN, including CIDR blocks, domain name, subnets, gateways, and creation date.
func describeVCN(v domain.VCN) string {
	parts := []string{}

	if len(v.CidrBlocks) > 0 {
		parts = append(parts, fmt.Sprintf("%d CIDRs", len(v.CidrBlocks)))
	}

	if v.DomainName != "" {
		parts = append(parts, v.DomainName)
	}

	if len(v.Subnets) > 0 {
		parts = append(parts, fmt.Sprintf("%d subnets", len(v.Subnets)))
	}
	if len(v.Gateways) > 0 {
		parts = append(parts, fmt.Sprintf("%d gateways", len(v.Gateways)))
	}

	if !v.TimeCreated.IsZero() {
		parts = append(parts, v.TimeCreated.Format("2006-01-02"))
	}

	return strings.Join(parts, " • ")
}
