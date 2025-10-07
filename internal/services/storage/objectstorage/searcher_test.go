package objectstorage

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearcher_GetFields(t *testing.T) {
	fields := GetSearchableFields()
	boost := GetBoostedFields()

	// basic expectations
	assert.Contains(t, fields, "Name")
	assert.Contains(t, fields, "OCID")
	assert.Contains(t, fields, "Namespace")
	assert.Contains(t, fields, "TagsKV")
	assert.Contains(t, fields, "TagsVal")

	// boosted are subset of fields
	for _, b := range boost {
		assert.Contains(t, fields, b)
	}
}

func TestSearchableBucket_ToIndexable_LowercasesAndMaps(t *testing.T) {
	b := Bucket{
		Name:               "Prod-Logs",
		OCID:               "OCID1.BUCKET.OC1..X",
		Namespace:          "MyNS",
		StorageTier:        "Standard",
		Visibility:         "Private",
		Encryption:         "SSE",
		Versioning:         "Enabled",
		ReplicationEnabled: true,
		IsReadOnly:         false,
		FreeformTags:       map[string]string{"Env": "Prod"},
	}
	m := SearchableBucket{b}.ToIndexable()

	// strings are lowercased
	for _, k := range []string{"Name", "OCID", "Namespace", "StorageTier", "Visibility", "Encryption", "Versioning", "TagsKV", "TagsVal"} {
		if v, ok := m[k].(string); ok {
			assert.Equal(t, strings.ToLower(v), v, k+" should be lowercased")
		}
	}
	// booleans converted to string
	assert.Equal(t, "true", m["ReplicationEnabled"])
	assert.Equal(t, "false", m["IsReadOnly"])
}
