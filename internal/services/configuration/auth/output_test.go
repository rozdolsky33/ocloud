package auth

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestGroupRegionsByPrefix(t *testing.T) {
	// Create test regions
	regions := []RegionInfo{
		{ID: "1", Name: "us-ashburn-1"},
		{ID: "2", Name: "us-phoenix-1"},
		{ID: "3", Name: "eu-frankfurt-1"},
		{ID: "4", Name: "ap-tokyo-1"},
		{ID: "5", Name: "uk-london-1"},
	}

	// Group the regions
	groups := groupRegionsByPrefix(regions)

	// Verify the groups
	assert.Len(t, groups, 4)
	assert.Len(t, groups["us"], 2)
	assert.Len(t, groups["eu"], 1)
	assert.Len(t, groups["ap"], 1)
	assert.Len(t, groups["uk"], 1)

	// Verify the contents of each group
	assert.Equal(t, "us-ashburn-1", groups["us"][0].Name)
	assert.Equal(t, "us-phoenix-1", groups["us"][1].Name)
	assert.Equal(t, "eu-frankfurt-1", groups["eu"][0].Name)
	assert.Equal(t, "ap-tokyo-1", groups["ap"][0].Name)
	assert.Equal(t, "uk-london-1", groups["uk"][0].Name)
}

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
		{"af", "Africa"},
		{"ca", "Canada"},
		{"il", "Israel"},
		{"me", "Middle East"},
		{"mx", "Mexico"},
		{"sa", "South America"},
		{"unknown", "unknown"}, // Unknown prefix should return itself
	}

	// Test each case
	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			title := getRegionGroupTitle(tc.prefix)
			assert.Equal(t, tc.expected, title)
		})
	}
}

func TestPrintExportVariable(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		tenancyName string
		compartment string
	}{
		{
			name:        "Both values provided",
			tenancyName: "example-tenancy",
			compartment: "example-compartment",
		},
		{
			name:        "Only tenancy name provided",
			tenancyName: "example-tenancy",
			compartment: "",
		},
		{
			name:        "Only compartment provided",
			tenancyName: "",
			compartment: "example-compartment",
		},
		{
			name:        "No values provided",
			tenancyName: "",
			compartment: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			err := PrintExportVariable(tc.tenancyName, tc.compartment)
			assert.NoError(t, err)

			// We can't easily test the output because it includes ANSI color codes
			// and is written directly to stdout. In a real-world scenario, we would
			// refactor the function to accept an io.Writer parameter to make it more testable.
			// For now, we just verify that the function doesn't return an error.
		})
	}
}

func TestDisplayRegionsTable(t *testing.T) {
	// Create test regions
	regions := []RegionInfo{
		{ID: "1", Name: "us-ashburn-1"},
		{ID: "2", Name: "us-phoenix-1"},
		{ID: "3", Name: "eu-frankfurt-1"},
	}

	// Test cases
	testCases := []struct {
		name   string
		filter string
		expect int // Number of region groups expected
	}{
		{
			name:   "No filter",
			filter: "",
			expect: 2, // us and eu groups
		},
		{
			name:   "Filter by us",
			filter: "us",
			expect: 1, // Only us group
		},
		{
			name:   "Filter by eu",
			filter: "eu",
			expect: 1, // Only eu group
		},
		{
			name:   "Filter by non-existent prefix",
			filter: "xx",
			expect: 0, // No groups
		},
	}

	// Test each case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer

			// Create a mock application context
			appCtx := &app.ApplicationContext{
				Stdout: &buf,
			}

			// Call the function
			err := DisplayRegionsTable(regions, appCtx, tc.filter)

			// Verify no error
			assert.NoError(t, err)

			// Get the output
			output := buf.String()

			// Count the number of region groups in the output
			// This is a simple heuristic - in a real test we might want to parse the output more carefully
			if tc.expect == 0 {
				assert.Empty(t, output)
			} else {
				assert.NotEmpty(t, output)
			}
		})
	}
}
