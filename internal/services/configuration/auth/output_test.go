package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestDisplayRegionsTable tests the DisplayRegionsTable function
func TestDisplayRegionsTable(t *testing.T) {
	// Create a sample regions list
	regions := []RegionInfo{
		{ID: "1", Name: "us-ashburn-1"},
		{ID: "2", Name: "us-phoenix-1"},
		{ID: "3", Name: "eu-frankfurt-1"},
	}

	// Test with no filter
	err := DisplayRegionsTable(regions, "")
	assert.NoError(t, err, "DisplayRegionsTable should not return an error with no filter")

	// Test with a filter
	err = DisplayRegionsTable(regions, "us")
	assert.NoError(t, err, "DisplayRegionsTable should not return an error with a filter")

	// Test with a filter that matches no regions
	err = DisplayRegionsTable(regions, "nonexistent")
	assert.NoError(t, err, "DisplayRegionsTable should not return an error with a non-matching filter")
}
