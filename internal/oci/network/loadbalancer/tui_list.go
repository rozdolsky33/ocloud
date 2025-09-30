package loadbalancer

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewLoadBalancerListModel builds a TUI list for load balancers.
func NewLoadBalancerListModel(lbs []domain.LoadBalancer) tui.Model {
	return tui.NewModel("Load Balancers", lbs, func(lb domain.LoadBalancer) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          lb.OCID,
			Title:       lb.Name,
			Description: description(lb),
		}
	})
}

func description(lb domain.LoadBalancer) string {
	ip := first(lb.IPAddresses)

	hs := healthSummary(lb.BackendHealth)

	line2 := joinNonEmpty(" • ",
		hs,
		firstNonEmpty(lb.VcnName, lb.VcnID),
	)

	return joinNonEmpty(" • ", ip, line2)
}

// --- helpers ---

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
	// consider any non-OK/non-UNKNOWN as unhealthy (CRITICAL, WARNING, etc.)
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
