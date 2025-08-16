package info

import (
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
}

// TestLoadTenancyMappings tests the LoadTenancyMappings method
func TestLoadTenancyMappings(t *testing.T) {
	// Create a new service
	service := NewService()

	// Test with a real implementation (skip if a file not available)
	t.Run("Real implementation", func(t *testing.T) {
		// Test with empty realm
		result, err := service.LoadTenancyMappings("")
		if err != nil {
			t.Skip("Skipping test because tenancy mapping file is not available")
		} else {
			assert.NotNil(t, result, "Result should not be nil")
			assert.NotNil(t, result.Mappings, "Mappings should not be nil")
		}

		// Test with a specific realm
		result, err = service.LoadTenancyMappings("OC1")
		if err != nil {
			t.Skip("Skipping test because tenancy mapping file is not available")
		} else {
			assert.NotNil(t, result, "Result should not be nil")
			assert.NotNil(t, result.Mappings, "Mappings should not be nil")
		}
	})
}
