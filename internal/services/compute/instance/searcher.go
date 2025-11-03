package instance

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableInstance is an adapter to make to compute.Instance searchable.
type SearchableInstance struct {
	compute.Instance
}

// ToIndexable converts an Instance to a map of searchable fields.
func (s SearchableInstance) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	// Join the security list and NSG names for searchability
	securityListNames := strings.Join(s.SecurityListNames, " ")
	nsgNames := strings.Join(s.NsgNames, " ")

	return map[string]any{
		"Name":          strings.ToLower(s.DisplayName),
		"Hostname":      strings.ToLower(s.Hostname),
		"PrimaryIP":     strings.ToLower(s.PrimaryIP),
		"ImageName":     strings.ToLower(s.ImageName),
		"ImageOS":       strings.ToLower(s.ImageOS),
		"Shape":         strings.ToLower(s.Shape),
		"OCID":          strings.ToLower(s.OCID),
		"FD":            strings.ToLower(s.FaultDomain),
		"AD":            strings.ToLower(s.AvailabilityDomain),
		"VcnName":       strings.ToLower(s.VcnName),
		"SubnetName":    strings.ToLower(s.SubnetName),
		"SecurityLists": strings.ToLower(securityListNames),
		"NSGs":          strings.ToLower(nsgNames),
		"TagsKV":        strings.ToLower(tagsKV),
		"TagsVal":       strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the list of fields to be indexed.
func GetSearchableFields() []string {
	return []string{
		"Name", "Hostname", "ImageName", "ImageOS", "Shape",
		"PrimaryIP", "OCID", "VcnName", "SubnetName", "FD", "AD",
		"SecurityLists", "NSGs",
		"TagsKV", "TagsVal",
	}
}

// GetBoostedFields returns the list of fields to be boosted in the search.
func GetBoostedFields() []string {
	return []string{"Name", "Hostname"}
}

// ToSearchableInstances converts a slice of compute.Instance to a slice of search.Indexable.
func ToSearchableInstances(instances []Instance) []search.Indexable {
	searchable := make([]search.Indexable, len(instances))
	for i, inst := range instances {
		searchable[i] = SearchableInstance{inst}
	}
	return searchable
}
