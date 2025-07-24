package info

import (
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestPrintMappingsFile tests the PrintMappingsFile function
func TestPrintMappingsFile(t *testing.T) {
	// This is a smoke test since we can't easily mock the printer.New function
	// Test with empty mappings and JSON output
	err := PrintMappingsFile([]appConfig.MappingsFile{}, true)
	assert.NoError(t, err, "PrintMappingsFile should not return an error with empty mappings and JSON output")

	// Test with empty mappings and table output
	err = PrintMappingsFile([]appConfig.MappingsFile{}, false)
	assert.NoError(t, err, "PrintMappingsFile should not return an error with empty mappings and table output")

	// Test with non-empty mappings and JSON output
	mappings := []appConfig.MappingsFile{
		{
			Realm:        "OC1",
			Environment:  "Test",
			Tenancy:      "TestTenancy",
			Compartments: "TestCompartment",
			Regions:      "us-ashburn-1",
		},
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

	err = PrintMappingsFile(mappings, true)
	assert.NoError(t, err, "PrintMappingsFile should not return an error with non-empty mappings and JSON output")

	// Test with non-empty mappings and table output
	// We're already redirecting stdout, so we can just call the function
	err = PrintMappingsFile(mappings, false)
	assert.NoError(t, err, "PrintMappingsFile should not return an error with non-empty mappings and table output")
}

// TestGroupMappingsByRealm tests the groupMappingsByRealm function
func TestGroupMappingsByRealm(t *testing.T) {
	// Test with empty mappings
	result := groupMappingsByRealm([]appConfig.MappingsFile{})
	assert.Empty(t, result, "groupMappingsByRealm should return an empty map for empty mappings")

	// Test with non-empty mappings
	mappings := []appConfig.MappingsFile{
		{
			Realm:        "OC1",
			Environment:  "Test1",
			Tenancy:      "TestTenancy1",
			Compartments: "TestCompartment1",
			Regions:      "us-ashburn-1",
		},
		{
			Realm:        "OC1",
			Environment:  "Test2",
			Tenancy:      "TestTenancy2",
			Compartments: "TestCompartment2",
			Regions:      "us-phoenix-1",
		},
		{
			Realm:        "OC2",
			Environment:  "Test3",
			Tenancy:      "TestTenancy3",
			Compartments: "TestCompartment3",
			Regions:      "eu-frankfurt-1",
		},
	}

	result = groupMappingsByRealm(mappings)
	assert.Len(t, result, 2, "groupMappingsByRealm should return a map with 2 entries")
	assert.Len(t, result["OC1"], 2, "groupMappingsByRealm should return a map with 2 entries for OC1")
	assert.Len(t, result["OC2"], 1, "groupMappingsByRealm should return a map with 1 entry for OC2")
}
