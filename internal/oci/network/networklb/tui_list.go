package networklb

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewNetworkLoadBalancerListModel builds a TUI list for network load balancers.
func NewNetworkLoadBalancerListModel(nlbs []domain.NetworkLoadBalancer) tui.Model {
	return tui.NewModel("Network Load Balancers", nlbs, func(nlb domain.NetworkLoadBalancer) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          nlb.OCID,
			Title:       nlb.Name,
			Description: nlbDescription(nlb),
		}
	})
}

func nlbDescription(nlb domain.NetworkLoadBalancer) string {
	ip := first(nlb.IPAddresses)
	hs := healthSummary(nlb.BackendHealth)
	line2 := joinNonEmpty(" • ",
		hs,
		firstNonEmpty(nlb.VcnName, nlb.VcnID),
	)
	return joinNonEmpty(" • ", ip, line2)
}

func first(ss []string) string {
	if len(ss) > 0 {
		return ss[0]
	}
	return ""
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func joinNonEmpty(sep string, parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return strings.Join(out, sep)
}

func healthSummary(health map[string]string) string {
	if len(health) == 0 {
		return "Health N/A"
	}
	total := len(health)
	ok := 0
	unknown := 0
	for _, s := range health {
		switch strings.ToUpper(s) {
		case "OK":
			ok++
		case "UNKNOWN":
			unknown++
		}
	}
	unhealthy := total - ok - unknown

	switch {
	case unhealthy > 0:
		return fmt.Sprintf("UNHEALTHY (%d/%d)", unhealthy, total)
	case ok == total:
		return fmt.Sprintf("Health OK (%d/%d)", ok, total)
	default:
		return fmt.Sprintf("Health %d OK, %d UNKNOWN", ok, unknown)
	}
}
