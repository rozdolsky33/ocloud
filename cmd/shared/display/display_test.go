package display

import (
	"os"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/stretchr/testify/assert"
)

// TestPrintOCIConfiguration tests the PrintOCIConfiguration function
// This is a simple test that verifies the function doesn't panic
// We can't easily test the actual output since it writes to stdout
func TestPrintOCIConfiguration(t *testing.T) {
	// Save original environment variables
	originalProfile := os.Getenv(flags.EnvKeyProfile)
	originalTenancy := os.Getenv(flags.EnvKeyTenancyName)
	originalCompartment := os.Getenv(flags.EnvKeyCompartment)

	// Restore environment variables after the test
	defer func() {
		os.Setenv(flags.EnvKeyProfile, originalProfile)
		os.Setenv(flags.EnvKeyTenancyName, originalTenancy)
		os.Setenv(flags.EnvKeyCompartment, originalCompartment)
	}()

	// Test case 1: No environment variables set
	os.Unsetenv(flags.EnvKeyProfile)
	os.Unsetenv(flags.EnvKeyTenancyName)
	os.Unsetenv(flags.EnvKeyCompartment)

	// This should not panic
	assert.NotPanics(t, func() {
		PrintOCIConfiguration()
	}, "PrintOCIConfiguration should not panic when no environment variables are set")

	// Test case 2: All environment variables set
	os.Setenv(flags.EnvKeyProfile, "test-profile")
	os.Setenv(flags.EnvKeyTenancyName, "test-tenancy")
	os.Setenv(flags.EnvKeyCompartment, "test-compartment")

	// This should not panic
	assert.NotPanics(t, func() {
		PrintOCIConfiguration()
	}, "PrintOCIConfiguration should not panic when all environment variables are set")

	// Test case 3: Some environment variables set
	os.Setenv(flags.EnvKeyProfile, "test-profile")
	os.Unsetenv(flags.EnvKeyTenancyName)
	os.Setenv(flags.EnvKeyCompartment, "test-compartment")

	// This should not panic
	assert.NotPanics(t, func() {
		PrintOCIConfiguration()
	}, "PrintOCIConfiguration should not panic when some environment variables are set")

	// Test case 4: Test with session validation functionality
	// This should not panic even if the oci command is not available
	assert.NotPanics(t, func() {
		PrintOCIConfiguration()
	}, "PrintOCIConfiguration should not panic when checking session validity")
}

// TestDisplayBanner indirectly tests the displayBanner function through PrintOCIConfiguration
// This is a simple test that verifies the function doesn't panic
func TestDisplayBanner(t *testing.T) {
	// This should not panic
	assert.NotPanics(t, func() {
		displayBanner()
	}, "displayBanner should not panic")
}
