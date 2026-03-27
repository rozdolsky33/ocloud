package networklb

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
)

// SearchableNetworkLoadBalancer adapts NetworkLoadBalancer to the search.Indexable interface.
type SearchableNetworkLoadBalancer struct {
	NetworkLoadBalancer
}

// ToIndexable converts a NetworkLoadBalancer to a map of searchable fields.
func (s SearchableNetworkLoadBalancer) ToIndexable() map[string]any {
	join := func(ss []string) string {
		if len(ss) == 0 {
			return ""
		}
		out := make([]string, 0, len(ss))
		for _, v := range ss {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			out = append(out, strings.ToLower(v))
		}
		return strings.Join(out, " ")
	}

	return map[string]any{
		"Name":        strings.ToLower(s.Name),
		"OCID":        strings.ToLower(s.OCID),
		"Type":        strings.ToLower(s.Type),
		"State":       strings.ToLower(s.State),
		"VcnName":     strings.ToLower(s.VcnName),
		"IPAddresses": join(s.IPAddresses),
		"Subnets":     join(s.Subnets),
	}
}

// GetSearchableFields returns the fields to index for network load balancers.
func GetSearchableFields() []string {
	return []string{"Name", "OCID", "Type", "State", "VcnName", "IPAddresses", "Subnets"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID"}
}

// ToSearchableNetworkLoadBalancers converts a slice of NetworkLoadBalancer to a slice of search.Indexable.
func ToSearchableNetworkLoadBalancers(items []NetworkLoadBalancer) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, it := range items {
		out[i] = SearchableNetworkLoadBalancer{it}
	}
	return out
}
