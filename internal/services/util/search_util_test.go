package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestFlattenTags tests the FlattenTags function
func TestFlattenTags(t *testing.T) {
	// Test with nil maps
	result, err := FlattenTags(nil, nil)
	assert.NoError(t, err, "FlattenTags should not return an error with nil maps")
	assert.Equal(t, "", result, "FlattenTags should return an empty string with nil maps")

	// Test with empty maps
	result, err = FlattenTags(map[string]string{}, map[string]map[string]interface{}{})
	assert.NoError(t, err, "FlattenTags should not return an error with empty maps")
	assert.Equal(t, "", result, "FlattenTags should return an empty string with empty maps")

	// Test with only freeform tags
	freeform := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	result, err = FlattenTags(freeform, nil)
	assert.NoError(t, err, "FlattenTags should not return an error with only freeform tags")
	assert.Contains(t, result, "key1:value1", "FlattenTags should include key1:value1")
	assert.Contains(t, result, "key2:value2", "FlattenTags should include key2:value2")

	// Test with only defined tags
	defined := map[string]map[string]interface{}{
		"namespace1": {
			"key1": "value1",
			"key2": 42,
		},
		"namespace2": {
			"key3": "value3",
		},
	}
	result, err = FlattenTags(nil, defined)
	assert.NoError(t, err, "FlattenTags should not return an error with only defined tags")
	assert.Contains(t, result, "namespace1.key1:value1", "FlattenTags should include namespace1.key1:value1")
	assert.Contains(t, result, "namespace1.key2:42", "FlattenTags should include namespace1.key2:42")
	assert.Contains(t, result, "namespace2.key3:value3", "FlattenTags should include namespace2.key3:value3")

	// Test with both freeform and defined tags
	result, err = FlattenTags(freeform, defined)
	assert.NoError(t, err, "FlattenTags should not return an error with both freeform and defined tags")
	assert.Contains(t, result, "key1:value1", "FlattenTags should include key1:value1")
	assert.Contains(t, result, "key2:value2", "FlattenTags should include key2:value2")
	assert.Contains(t, result, "namespace1.key1:value1", "FlattenTags should include namespace1.key1:value1")
	assert.Contains(t, result, "namespace1.key2:42", "FlattenTags should include namespace1.key2:42")
	assert.Contains(t, result, "namespace2.key3:value3", "FlattenTags should include namespace2.key3:value3")

	// Test with empty keys and values
	freeform = map[string]string{
		"":     "value",
		"key":  "",
		"key3": "value3",
	}
	defined = map[string]map[string]interface{}{
		"": {
			"key": "value",
		},
		"namespace": {
			"":    "value",
			"key": nil,
		},
	}
	result, err = FlattenTags(freeform, defined)
	assert.NoError(t, err, "FlattenTags should not return an error with empty keys and values")
	assert.NotContains(t, result, ":value\n", "FlattenTags should skip empty keys")
	assert.NotContains(t, result, "key:\n", "FlattenTags should skip empty values")
	assert.Contains(t, result, "key3:value3", "FlattenTags should include valid key-value pairs")
}

// TestExtractTagValues tests the ExtractTagValues function
func TestExtractTagValues(t *testing.T) {
	// Test with nil maps
	result, err := ExtractTagValues(nil, nil)
	assert.NoError(t, err, "ExtractTagValues should not return an error with nil maps")
	assert.Equal(t, "", result, "ExtractTagValues should return an empty string with nil maps")

	// Test with empty maps
	result, err = ExtractTagValues(map[string]string{}, map[string]map[string]interface{}{})
	assert.NoError(t, err, "ExtractTagValues should not return an error with empty maps")
	assert.Equal(t, "", result, "ExtractTagValues should return an empty string with empty maps")

	// Test with only freeform tags
	freeform := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	result, err = ExtractTagValues(freeform, nil)
	assert.NoError(t, err, "ExtractTagValues should not return an error with only freeform tags")
	assert.Contains(t, result, "value1", "ExtractTagValues should include value1")
	assert.Contains(t, result, "value2", "ExtractTagValues should include value2")

	// Test with only defined tags
	defined := map[string]map[string]interface{}{
		"namespace1": {
			"key1": "value1",
			"key2": 42,
		},
		"namespace2": {
			"key3": "value3",
		},
	}
	result, err = ExtractTagValues(nil, defined)
	assert.NoError(t, err, "ExtractTagValues should not return an error with only defined tags")
	assert.Contains(t, result, "value1", "ExtractTagValues should include value1")
	assert.Contains(t, result, "42", "ExtractTagValues should include 42")
	assert.Contains(t, result, "value3", "ExtractTagValues should include value3")

	// Test with both freeform and defined tags
	result, err = ExtractTagValues(freeform, defined)
	assert.NoError(t, err, "ExtractTagValues should not return an error with both freeform and defined tags")
	assert.Contains(t, result, "value1", "ExtractTagValues should include value1")
	assert.Contains(t, result, "value2", "ExtractTagValues should include value2")
	assert.Contains(t, result, "42", "ExtractTagValues should include 42")
	assert.Contains(t, result, "value3", "ExtractTagValues should include value3")

	// Test with empty values
	freeform = map[string]string{
		"key1": "",
		"key2": "value2",
	}
	defined = map[string]map[string]interface{}{
		"namespace1": {
			"key1": nil,
			"key2": "",
		},
		"namespace2": {
			"key3": "value3",
		},
	}
	result, err = ExtractTagValues(freeform, defined)
	assert.NoError(t, err, "ExtractTagValues should not return an error with empty values")
	// We can't use NotContains with an empty string because it would match the spaces between values
	// Instead, we'll check that the result doesn't contain empty values by ensuring it doesn't have consecutive spaces
	assert.NotContains(t, result, "  ", "ExtractTagValues should not have consecutive spaces (empty values)")
	assert.Contains(t, result, "value2", "ExtractTagValues should include valid values")
	assert.Contains(t, result, "value3", "ExtractTagValues should include valid values")
}
