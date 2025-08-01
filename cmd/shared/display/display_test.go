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
	originalProfile := os.Getenv("OCI_CLI_PROFILE")
	originalTenancy := os.Getenv(flags.EnvOCITenancyName)
	originalCompartment := os.Getenv(flags.EnvOCICompartment)

	// Restore environment variables after the test
	defer func() {
		os.Setenv("OCI_CLI_PROFILE", originalProfile)
		os.Setenv(flags.EnvOCITenancyName, originalTenancy)
		os.Setenv(flags.EnvOCICompartment, originalCompartment)
	}()

	// Test case 1: No environment variables set
	os.Unsetenv("OCI_CLI_PROFILE")
	os.Unsetenv(flags.EnvOCITenancyName)
	os.Unsetenv(flags.EnvOCICompartment)

	// This should not panic
	assert.NotPanics(t, func() {
		PrintOCIConfiguration()
	}, "PrintOCIConfiguration should not panic when no environment variables are set")

	// Test case 2: All environment variables set
	os.Setenv("OCI_CLI_PROFILE", "test-profile")
	os.Setenv(flags.EnvOCITenancyName, "test-tenancy")
	os.Setenv(flags.EnvOCICompartment, "test-compartment")

	// This should not panic
	assert.NotPanics(t, func() {
		PrintOCIConfiguration()
	}, "PrintOCIConfiguration should not panic when all environment variables are set")

	// Test case 3: Some environment variables set
	os.Setenv("OCI_CLI_PROFILE", "test-profile")
	os.Unsetenv(flags.EnvOCITenancyName)
	os.Setenv(flags.EnvOCICompartment, "test-compartment")

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
