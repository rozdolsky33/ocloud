package policy

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchablePolicy adapts identity.Policy to the search.Indexable interface.
type SearchablePolicy struct {
	Policy
}

// ToIndexable converts a Policy to a map of searchable fields.
func (s SearchablePolicy) ToIndexable() map[string]any {
	// Flatten tags in two flavors: key:value form and just values
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)
	// Join statements into a single string
	stmt := strings.ToLower(strings.Join(s.Statement, " "))

	return map[string]any{
		"Name":        strings.ToLower(s.Name),
		"Description": strings.ToLower(s.Description),
		"OCID":        strings.ToLower(s.ID),
		"Statements":  stmt,
		"TagsKV":      strings.ToLower(tagsKV),
		"TagsVal":     strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the fields to index for policies.
func GetSearchableFields() []string {
	return []string{"Name", "Description", "OCID", "Statements", "TagsKV", "TagsVal"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "Statements", "TagsKV", "TagsVal"}
}

// ToSearchablePolicies converts a slice of Policy to a slice of search.Indexable.
func ToSearchablePolicies(items []Policy) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, p := range items {
		out[i] = SearchablePolicy{p}
	}
	return out
}
