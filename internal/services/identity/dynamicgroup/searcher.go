package dynamicgroup

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableDynamicGroup adapts identity.DynamicGroup to the search.Indexable interface.
type SearchableDynamicGroup struct {
	DynamicGroup
}

// ToIndexable converts a DynamicGroup to a map of searchable fields.
func (s SearchableDynamicGroup) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	return map[string]any{
		"Name":         strings.ToLower(s.Name),
		"Description":  strings.ToLower(s.Description),
		"OCID":         strings.ToLower(s.OCID),
		"State":        strings.ToLower(s.LifecycleState),
		"MatchingRule": strings.ToLower(s.MatchingRule),
		"TagsKV":       strings.ToLower(tagsKV),
		"TagsVal":      strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the fields to index for dynamic groups.
func GetSearchableFields() []string {
	return []string{"Name", "Description", "OCID", "State", "MatchingRule", "TagsKV", "TagsVal"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "TagsKV", "TagsVal"}
}

// ToSearchableDynamicGroups converts a slice of DynamicGroup to a slice of search.Indexable.
func ToSearchableDynamicGroups(items []DynamicGroup) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, c := range items {
		out[i] = SearchableDynamicGroup{c}
	}
	return out
}
