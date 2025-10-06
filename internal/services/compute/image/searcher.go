package image

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/services/search"
)

// SearchableImage is an adapter to make to compute.Image searchable.
type SearchableImage struct {
	compute.Image
}

// ToIndexable converts an Image to a map of searchable fields.
func (s SearchableImage) ToIndexable() map[string]any {
	return map[string]any{
		"Name":            strings.ToLower(s.DisplayName),
		"Created":         strings.ToLower(s.TimeCreated.Format("2006-01-02 15:04:05")),
		"OperatingSystem": strings.ToLower(s.OperatingSystem),
		"OSVersion":       strings.ToLower(s.OperatingSystemVersion),
		"OCID":            strings.ToLower(s.OCID),
		"LaunchMode":      strings.ToLower(s.LaunchMode),
	}
}

// GetSearchableFields returns the list of fields to be indexed.
func GetSearchableFields() []string {
	return []string{"Name", "OperatingSystem", "OSVersion", "OCID", "LaunchMode"}
}

// GetBoostedFields returns the list of fields to be boosted in the search.
func GetBoostedFields() []string {
	return []string{"Name", "OperatingSystem", "OSVersion"}
}

// ToSearchableImages converts a slice of compute.Image to a slice of search.Indexable.
func ToSearchableImages(images []Image) []search.Indexable {
	searchable := make([]search.Indexable, len(images))
	for i, img := range images {
		searchable[i] = SearchableImage{img}
	}
	return searchable
}
