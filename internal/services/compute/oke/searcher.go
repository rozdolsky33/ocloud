package oke

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableCluster adapts Cluster to the search.Indexable interface.
type SearchableCluster struct {
	Cluster
}

// ToIndexable converts a Cluster to a map of searchable fields.
func (s SearchableCluster) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	npNames := make([]string, 0, len(s.NodePools))
	npShapes := make([]string, 0, len(s.NodePools))
	for _, np := range s.NodePools {
		npNames = append(npNames, np.DisplayName)
		npShapes = append(npShapes, np.NodeShape)
	}

	return map[string]any{
		"Name":         strings.ToLower(s.DisplayName),
		"OCID":         strings.ToLower(s.OCID),
		"K8sVersion":   strings.ToLower(s.KubernetesVersion),
		"State":        strings.ToLower(s.State),
		"VcnOCID":      strings.ToLower(s.VcnOCID),
		"PrivEndpoint": strings.ToLower(s.PrivateEndpoint),
		"PubEndpoint":  strings.ToLower(s.PublicEndpoint),
		"NodePools":    strings.ToLower(strings.Join(npNames, ",")),
		"NodeShapes":   strings.ToLower(strings.Join(npShapes, ",")),
		"TagsKV":       strings.ToLower(tagsKV),
		"TagsVal":      strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the list of fields to be indexed for OKE clusters.
func GetSearchableFields() []string {
	return []string{
		"Name", "OCID", "K8sVersion", "State", "VcnOCID",
		"PrivEndpoint", "PubEndpoint", "NodePools", "NodeShapes", "TagsKV", "TagsVal",
	}
}

// GetBoostedFields returns the list of fields to be boosted in the search.
func GetBoostedFields() []string {
	return []string{"Name", "NodePools", "NodeShapes"}
}

// ToSearchableClusters converts a slice of Cluster to a slice of search.Indexable.
func ToSearchableClusters(clusters []Cluster) []search.Indexable {
	searchable := make([]search.Indexable, len(clusters))
	for i, c := range clusters {
		searchable[i] = SearchableCluster{c}
	}
	return searchable
}
