package resources

import (
	"context"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// Setup test environment
func setupTest(t *testing.T) {
	// Initialize logger for tests
	logger.InitLogger(logger.CmdLogger)
}

// TestListInstances tests the ListInstances function
func TestListInstances(t *testing.T) {
	setupTest(t)

	// Create a mock AppContext
	mockCtx := &app.AppContext{
		Ctx:             context.Background(),
		Provider:        common.DefaultConfigProvider(),
		CompartmentID:   "mock-compartment-id",
		CompartmentName: "mock-compartment",
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

	// Create a mock AppContext
	mockCtx := &app.AppContext{
		Ctx:             context.Background(),
		Provider:        common.DefaultConfigProvider(),
		CompartmentID:   "mock-compartment-id",
		CompartmentName: "mock-compartment",
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
