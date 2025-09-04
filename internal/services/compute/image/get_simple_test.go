package image

import (
	"io"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestGetImagesSimple is a simplified test for the GetImages function
// that doesn't rely on mocking the OCI SDK interfaces
func TestGetImagesSimple(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for ListImages since it requires the OCI SDK")

	// In a real test, we would:
	// 1. Create a mock application context
	// 2. Create mock images
	// 3. Call ListImages with different parameters
	// 4. Verify the results

	appCtx := &app.ApplicationContext{
		TenancyName:     "TestTenancy",
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
		Stdout:          io.Discard, // Discard output to avoid cluttering the test output
	}

	err := GetImages(appCtx, 20, 1, false)

	assert.NoError(t, err)
}
