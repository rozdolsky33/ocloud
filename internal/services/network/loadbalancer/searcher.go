package loadbalancer

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
)

// SearchableLoadBalancer adapts LoadBalancer to the search.Indexable interface.
type SearchableLoadBalancer struct {
	LoadBalancer
}

// ToIndexable converts a LoadBalancer to a map of searchable fields.
func (s SearchableLoadBalancer) ToIndexable() map[string]any {
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
		"Name":            strings.ToLower(s.Name),
		"OCID":            strings.ToLower(s.OCID),
		"Type":            strings.ToLower(s.Type),
		"State":           strings.ToLower(s.State),
		"VcnName":         strings.ToLower(s.VcnName),
		"Shape":           strings.ToLower(s.Shape),
		"IPAddresses":     join(s.IPAddresses),
		"Hostnames":       join(s.Hostnames),
		"SSLCertificates": join(s.SSLCertificates),
		"Subnets":         join(s.Subnets),
	}
}

// GetSearchableFields returns the fields to index for load balancers.
func GetSearchableFields() []string {
	return []string{"Name", "OCID", "Type", "State", "VcnName", "Shape", "IPAddresses", "Hostnames", "SSLCertificates", "Subnets"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "Hostnames"}
}

// ToSearchableLoadBalancers converts a slice of LoadBalancer to a slice of search.Indexable.
func ToSearchableLoadBalancers(items []LoadBalancer) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, it := range items {
		out[i] = SearchableLoadBalancer{it}
	}
	return out
}
