package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Setup test environment
func setupTest(t *testing.T) {
	// Initialize logger for tests
	logger.InitLogger(logger.CmdLogger)
}

// TestListInstances tests the ListInstances function
func TestListInstances(t *testing.T) {
	setupTest(t)

	// Create a mock AppContext with our mock provider
	mockCtx := &app.AppContext{
		Provider:        oci.NewMockConfigurationProvider(),
		CompartmentID:   "mock-compartment-id",
		CompartmentName: "mock-compartment",
		Logger:          logger.CmdLogger,
	}

	// Test successful case
	err := ListInstances(mockCtx)
	assert.NoError(t, err)

	// Test error case (can't easily test this without mocking the oci package)
	// This would require a more complex test setup with dependency injection
}

// TestFindInstances tests the FindInstances function
func TestFindInstances(t *testing.T) {
	setupTest(t)

	// Create a mock AppContext with our mock provider
	mockCtx := &app.AppContext{
		Provider:        oci.NewMockConfigurationProvider(),
		CompartmentID:   "mock-compartment-id",
		CompartmentName: "mock-compartment",
		Logger:          logger.CmdLogger,
	}

	// Test successful case
	err := FindInstances(mockCtx, "test-pattern", false)
	assert.NoError(t, err)

	// Test with image details
	err = FindInstances(mockCtx, "test-pattern", true)
	assert.NoError(t, err)

	// Test error case (can't easily test this without mocking the oci package)
	// This would require a more complex test setup with dependency injection
}
