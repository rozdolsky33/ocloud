package compartment

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableCompartment adapts identity.Compartment to the search.Indexable interface.
type SearchableCompartment struct {
	Compartment
}

// ToIndexable converts a Compartment to a map of searchable fields.
func (s SearchableCompartment) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	return map[string]any{
		"Name":        strings.ToLower(s.DisplayName),
		"Description": strings.ToLower(s.Description),
		"OCID":        strings.ToLower(s.OCID),
		"State":       strings.ToLower(s.LifecycleState),
		"TagsKV":      strings.ToLower(tagsKV),
		"TagsVal":     strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the fields to index for compartments.
func GetSearchableFields() []string {
	return []string{"Name", "Description", "OCID", "State", "TagsKV", "TagsVal"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "TagsKV", "TagsVal"}
}

// ToSearchableCompartments converts a slice of Compartment to a slice of search.Indexable.
func ToSearchableCompartments(items []Compartment) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, c := range items {
		out[i] = SearchableCompartment{c}
	}
	return out
}
