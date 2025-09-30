package loadbalancer

import (
	"fmt"
	"sort"
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
	meta := joinNonEmpty(" • ",
		lb.Type,
		ip,
	)

	// 2) Listeners summary: "2 listeners (https:443 → default)"
	ls := listenerSummary(lb.Listeners)

	// 3) Health summary: "Health OK (2/2)" or "UNHEALTHY (1/3)"
	hs := healthSummary(lb.BackendHealth)

	line2 := joinNonEmpty(" • ",
		ls,
		hs,
		firstNonEmpty(lb.VcnName, lb.VcnID),
	)

	return joinNonEmpty(" • ", meta, line2)
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

func listenerSummary(listeners map[string]string) string {
	if len(listeners) == 0 {
		return "0 listeners"
	}
	names := make([]string, 0, len(listeners))
	for n := range listeners {
		names = append(names, n)
	}
	sort.Strings(names)
	example := listeners[names[0]]
	if len(example) > 40 {
		example = example[:40] + "…"
	}
	if len(listeners) == 1 {
		return "1 listener (" + example + ")"
	}
	return fmt.Sprintf("%d listeners (%s)", len(listeners), example)
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
