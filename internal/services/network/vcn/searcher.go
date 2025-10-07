package vcn

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableVCN adapts VCN to the search.Indexable interface.
type SearchableVCN struct {
	VCN
}

// ToIndexable converts a VCN to a map of searchable fields.
func (s SearchableVCN) ToIndexable() map[string]any {
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

	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	// Collect names of related resources for better search coverage.
	gwNames := make([]string, 0, len(s.Gateways))
	for _, g := range s.Gateways {
		gwNames = append(gwNames, g.DisplayName)
	}
	subnetNames := make([]string, 0, len(s.Subnets))
	for _, sn := range s.Subnets {
		subnetNames = append(subnetNames, sn.DisplayName)
	}
	nsgNames := make([]string, 0, len(s.NSGs))
	for _, n := range s.NSGs {
		nsgNames = append(nsgNames, n.DisplayName)
	}
	rtNames := make([]string, 0, len(s.RouteTables))
	for _, r := range s.RouteTables {
		rtNames = append(rtNames, r.DisplayName)
	}
	slNames := make([]string, 0, len(s.SecurityLists))
	for _, sl := range s.SecurityLists {
		slNames = append(slNames, sl.DisplayName)
	}

	return map[string]any{
		"Name":        strings.ToLower(s.DisplayName),
		"OCID":        strings.ToLower(s.OCID),
		"State":       strings.ToLower(s.LifecycleState),
		"CIDRs":       join(s.CidrBlocks),
		"DnsLabel":    strings.ToLower(s.DnsLabel),
		"DomainName":  strings.ToLower(s.DomainName),
		"TagsKV":      strings.ToLower(tagsKV),
		"TagsVal":     strings.ToLower(tagsVal),
		"Gateways":    join(gwNames),
		"Subnets":     join(subnetNames),
		"NSGs":        join(nsgNames),
		"RouteTables": join(rtNames),
		"SecLists":    join(slNames),
	}
}

// GetSearchableFields returns the fields to index for VCNs.
func GetSearchableFields() []string {
	return []string{"Name", "OCID", "State", "CIDRs", "DnsLabel", "DomainName", "TagsKV", "TagsVal", "Gateways", "Subnets", "NSGs", "RouteTables", "SecLists"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "DnsLabel", "DomainName", "TagsKV", "TagsVal"}
}

// ToSearchableVCNs converts a slice of VCN to a slice of search.Indexable.
func ToSearchableVCNs(items []VCN) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, it := range items {
		out[i] = SearchableVCN{it}
	}
	return out
}
