package objectstorage

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableBucket adapts Bucket to the search.Indexable interface.
type SearchableBucket struct {
	Bucket
}

// ToIndexable converts a Bucket to a map of searchable fields.
func (s SearchableBucket) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	boolStr := func(b bool) string {
		if b {
			return "true"
		}
		return "false"
	}

	return map[string]any{
		"Name":               strings.ToLower(s.Name),
		"OCID":               strings.ToLower(s.OCID),
		"Namespace":          strings.ToLower(s.Namespace),
		"StorageTier":        strings.ToLower(s.StorageTier),
		"Visibility":         strings.ToLower(s.Visibility),
		"Encryption":         strings.ToLower(s.Encryption),
		"Versioning":         strings.ToLower(s.Versioning),
		"ReplicationEnabled": boolStr(s.ReplicationEnabled),
		"IsReadOnly":         boolStr(s.IsReadOnly),
		"TagsKV":             strings.ToLower(tagsKV),
		"TagsVal":            strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the fields to index for Buckets.
func GetSearchableFields() []string {
	return []string{"Name", "OCID", "Namespace", "StorageTier", "Visibility", "Encryption", "Versioning", "ReplicationEnabled", "IsReadOnly", "TagsKV", "TagsVal"}
}

// GetBoostedFields returns fields to boost during the search for better relevance.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "Namespace", "TagsKV", "TagsVal"}
}

// ToSearchableBuckets converts a slice of Bucket to a slice of search.Indexable.
func ToSearchableBuckets(items []Bucket) []search.Indexable {
	out := make([]search.Indexable, len(items))
	for i, it := range items {
		out[i] = SearchableBucket{it}
	}
	return out
}
