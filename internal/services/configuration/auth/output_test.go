package auth

import (
	"github.com/stretchr/testify/assert"
	"os"
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

	// Save the original stdout
	oldStdout := os.Stdout

	// Create a temporary file to redirect stdout
	tmpFile, err := os.CreateTemp("", "output_test")
	assert.NoError(t, err, "Failed to create temporary file")
	defer os.Remove(tmpFile.Name())

	// Redirect stdout to the temporary file
	os.Stdout = tmpFile
	defer func() {
		os.Stdout = oldStdout
	}()

	// Test with no filter
	err = DisplayRegionsTable(regions, "")
	assert.NoError(t, err, "DisplayRegionsTable should not return an error with no filter")

	// Test with a filter
	err = DisplayRegionsTable(regions, "us")
	assert.NoError(t, err, "DisplayRegionsTable should not return an error with a filter")

	// Test with a filter that matches no regions
	err = DisplayRegionsTable(regions, "nonexistent")
	assert.NoError(t, err, "DisplayRegionsTable should not return an error with a non-matching filter")
}

// TestGroupRegionsByPrefix tests the groupRegionsByPrefix function
func TestGroupRegionsByPrefix(t *testing.T) {
	// Create a sample regions list
	regions := []RegionInfo{
		{ID: "1", Name: "us-ashburn-1"},
		{ID: "2", Name: "us-phoenix-1"},
		{ID: "3", Name: "eu-frankfurt-1"},
		{ID: "4", Name: "ap-tokyo-1"},
		{ID: "5", Name: "ap-seoul-1"},
	}

	// Group regions by prefix
	result := groupRegionsByPrefix(regions)

	// Verify the result
	assert.Len(t, result, 3, "Should have 3 region groups (us, eu, ap)")
	assert.Len(t, result["us"], 2, "Should have 2 US regions")
	assert.Len(t, result["eu"], 1, "Should have 1 EU region")
	assert.Len(t, result["ap"], 2, "Should have 2 AP regions")

	// Verify the contents of each group
	assert.Equal(t, "us-ashburn-1", result["us"][0].Name, "First US region should be us-ashburn-1")
	assert.Equal(t, "us-phoenix-1", result["us"][1].Name, "Second US region should be us-phoenix-1")
	assert.Equal(t, "eu-frankfurt-1", result["eu"][0].Name, "EU region should be eu-frankfurt-1")
	assert.Equal(t, "ap-tokyo-1", result["ap"][0].Name, "First AP region should be ap-tokyo-1")
	assert.Equal(t, "ap-seoul-1", result["ap"][1].Name, "Second AP region should be ap-seoul-1")

	// Test with empty regions
	emptyResult := groupRegionsByPrefix([]RegionInfo{})
	assert.Len(t, emptyResult, 0, "Should have 0 region groups for empty input")

	// Test with invalid region name format
	invalidRegions := []RegionInfo{
		{ID: "1", Name: "invalid"},
	}
	invalidResult := groupRegionsByPrefix(invalidRegions)
	assert.Len(t, invalidResult, 1, "Should have 1 region group for invalid input")
	assert.Len(t, invalidResult["invalid"], 1, "Should have 1 region in the invalid group")
}

// TestGetRegionGroupTitle tests the getRegionGroupTitle function
func TestGetRegionGroupTitle(t *testing.T) {
	// Test cases
	testCases := []struct {
		prefix   string
		expected string
	}{
		{"us", "United States"},
		{"eu", "Europe"},
		{"ap", "Asia Pacific"},
		{"uk", "United Kingdom"},
		{"ca", "Canada"},
		{"sa", "South America"},
		{"me", "Middle East"},
		{"af", "Africa"},
		{"il", "Israel"},
		{"mx", "Mexico"},
		{"unknown", "unknown"}, // Unknown prefix should return the prefix itself
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			result := getRegionGroupTitle(tc.prefix)
			assert.Equal(t, tc.expected, result, "getRegionGroupTitle(%s) should return %s", tc.prefix, tc.expected)
		})
	}
}
