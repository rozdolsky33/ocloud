package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceTags(t *testing.T) {
	t.Run("new resource tags", func(t *testing.T) {
		tags := ResourceTags{
			"environment": "production",
			"team":        "backend",
			"cost-center": "engineering",
		}

		assert.NotNil(t, tags)
		assert.Equal(t, "production", tags["environment"])
		assert.Equal(t, "backend", tags["team"])
		assert.Equal(t, "engineering", tags["cost-center"])
		assert.Len(t, tags, 3)
	})

	t.Run("empty resource tags", func(t *testing.T) {
		tags := ResourceTags{}
		assert.NotNil(t, tags)
		assert.Len(t, tags, 0)
	})

	t.Run("nil resource tags", func(t *testing.T) {
		var tags ResourceTags
		assert.Nil(t, tags)
		assert.Len(t, tags, 0)
	})

	t.Run("modify resource tags", func(t *testing.T) {
		tags := make(ResourceTags)
		tags["key1"] = "value1"
		tags["key2"] = "value2"

		assert.Len(t, tags, 2)
		assert.Equal(t, "value1", tags["key1"])
		assert.Equal(t, "value2", tags["key2"])

		// Update the existing key
		tags["key1"] = "updated-value1"
		assert.Equal(t, "updated-value1", tags["key1"])
		assert.Len(t, tags, 2)

		// Delete key
		delete(tags, "key2")
		assert.Len(t, tags, 1)
		assert.Equal(t, "", tags["key2"]) // Zero value for a missing key
	})
}
