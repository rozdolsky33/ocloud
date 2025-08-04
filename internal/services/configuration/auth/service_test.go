package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewService tests the NewService function
func TestNewService(t *testing.T) {
	// Create a new service
	service := NewService()

	// Verify that the service is not nil
	assert.NotNil(t, service, "Service should not be nil")
	assert.NotNil(t, service.logger, "Logger should not be nil")
	assert.NotNil(t, service.Provider, "Provider should not be nil")
}

// TestGetOCIRegions tests the getOCIRegions method
func TestGetOCIRegions(t *testing.T) {
	// Create a new service
	service := NewService()

	// Get the OCI regions
	regions := service.getOCIRegions()

	// Verify that the region list is not empty
	assert.NotEmpty(t, regions, "Regions list should not be empty")

	// Verify that the region list contains expected regions
	expectedRegions := []string{
		"us-ashburn-1",
		"us-phoenix-1",
		"eu-frankfurt-1",
		"ap-tokyo-1",
	}

	// Check that each expected region is in the list
	for _, expected := range expectedRegions {
		found := false
		for _, region := range regions {
			if region.Name == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected region %s not found in regions list", expected)
	}

	// Verify that each region has a non-empty ID and Name
	for _, region := range regions {
		assert.NotEmpty(t, region.ID, "Region ID should not be empty")
		assert.NotEmpty(t, region.Name, "Region Name should not be empty")
	}
}

// TestViewConfigurationWithErrorHandling tests the viewConfigurationWithErrorHandling method
// This is a limited test since it depends on the external state
func TestViewConfigurationWithErrorHandling(t *testing.T) {
	// Skip this test in normal test runs since it depends on the external state
	t.Skip("Skipping test that depends on external state")

	// In a real test environment, we would mock the info.ViewConfiguration function
	// For example:
	// mockViewConfiguration := func(useJSON bool, realm string) error {
	//     if realm == "nonexistent" {
	//         return fmt.Errorf("tenancy mapping file not found")
	//     }
	//     return nil
	// }
	//
	// // Save the original function and restore it after the test
	// originalViewConfiguration := info.ViewConfiguration
	// info.ViewConfiguration = mockViewConfiguration
	// defer func() { info.ViewConfiguration = originalViewConfiguration }()
	//
	// service := NewService()
	//
	// // Test with a realm that causes a "tenancy mapping file not found" error
	// err := service.viewConfigurationWithErrorHandling("nonexistent")
	// assert.NoError(t, err, "Should handle 'tenancy mapping file not found' error")
	//
	// // Test with a realm that doesn't cause an error
	// err = service.viewConfigurationWithErrorHandling("OC1")
	// assert.NoError(t, err, "Should not return an error for valid realm")
}

// TestPrintExportVariable tests the PrintExportVariable function
func TestPrintExportVariable(t *testing.T) {
	// Save the original stdout
	oldStdout := os.Stdout

	// Create a temporary file to redirect stdout
	tmpFile, err := os.CreateTemp("", "export_test")
	assert.NoError(t, err, "Failed to create temporary file")
	defer os.Remove(tmpFile.Name())

	// Redirect stdout to the temporary file
	os.Stdout = tmpFile
	defer func() {
		os.Stdout = oldStdout
	}()

	// Test with all parameters provided
	err = PrintExportVariable("DEFAULT", "test-tenancy", "test-compartment")
	assert.NoError(t, err, "PrintExportVariable should not return an error with all parameters")

	// Test with some parameters empty
	err = PrintExportVariable("DEFAULT", "", "")
	assert.NoError(t, err, "PrintExportVariable should not return an error with some empty parameters")

	// Test with all parameters empty
	err = PrintExportVariable("", "", "")
	assert.NoError(t, err, "PrintExportVariable should not return an error with all empty parameters")
}
