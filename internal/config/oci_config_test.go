package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/config/flags"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTest prepares the test environment and returns a cleanup function
func setupTest(t *testing.T) func() {
	// Save original environment variables to restore later
	originalEnvVars := map[string]string{
		flags.EnvKeyProfile:        os.Getenv(flags.EnvKeyProfile),
		flags.EnvKeyTenancyMapPath: os.Getenv(flags.EnvKeyTenancyMapPath),
	}

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "oci-config-test")
	require.NoError(t, err)

	// Return a cleanup function
	return func() {
		// Restore original environment variables
		for key, value := range originalEnvVars {
			if value == "" {
				err := os.Unsetenv(key)
				if err != nil {
					t.Error(err)
					return
				}
			} else {
				err := os.Setenv(key, value)
				if err != nil {
					t.Error(err)
					return
				}
			}
		}
		// Remove the temporary directory
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

// TestGetOCIProfile tests the GetOCIProfile function
func TestGetOCIProfile(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Test default profile when environment variable is not set
	err := os.Unsetenv(flags.EnvKeyProfile)
	if err != nil {
		t.Error(err)
		return
	}
	profile := GetOCIProfile()
	assert.Equal(t, flags.DefaultProfileName, profile)

	// Test custom profile when environment variable is set
	customProfile := "CUSTOM_PROFILE"
	err = os.Setenv(flags.EnvKeyProfile, customProfile)
	if err != nil {
		t.Error(err)
		return
	}
	profile = GetOCIProfile()
	assert.Equal(t, customProfile, profile)
}

// TestTenancyMapPath tests the TenancyMapPath function
func TestTenancyMapPath(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Test default path when environment variable is not set
	err := os.Unsetenv(flags.EnvKeyTenancyMapPath)
	if err != nil {
		t.Error(err)
		return
	}

	path := TenancyMapPath()
	assert.Equal(t, DefaultTenancyMapPath, path)

	// Test custom path when environment variable is set
	customPath := "/custom/path/to/tenancy-map.yaml"
	err = os.Setenv(flags.EnvKeyTenancyMapPath, customPath)
	if err != nil {
		return
	}
	path = TenancyMapPath()
	assert.Equal(t, customPath, path)
}

// TestEnsureFile tests the ensureFile function
func TestEnsureFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "oci-config-test")
	require.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Error(err)
		}
	}(tempDir)

	// Test with a file that doesn't exist
	nonExistentFile := filepath.Join(tempDir, "non-existent.txt")
	err = ensureFile(nonExistentFile)
	assert.Error(t, err)

	// Test with a file that exists
	existingFile := filepath.Join(tempDir, "existing.txt")
	err = os.WriteFile(existingFile, []byte("test"), 0644)
	require.NoError(t, err)
	err = ensureFile(existingFile)
	assert.NoError(t, err)

	// Test with a directory
	err = ensureFile(tempDir)
	assert.Error(t, err)
}

// TestLoadTenancyMap tests the LoadTenancyMap function
func TestLoadTenancyMap(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "oci-config-test")
	require.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Error(err)
		}
	}(tempDir)

	// Create a valid tenancy map file
	validMapContent := `
- environment: test
  tenancy: test-tenancy
  tenancy_id: ocid1.tenancy.oc1..test
  realm: test-realm
  compartments: 
    - test-compartment-1
    - test-compartment-2
  regions: 
    - test-region-1
    - test-region-2
`
	var validMapFile = filepath.Join(tempDir, "valid-tenancy-map.yaml")
	err = os.WriteFile(validMapFile, []byte(validMapContent), 0644)
	require.NoError(t, err)

	// Set the environment variable to point to the test file
	err = os.Setenv(flags.EnvKeyTenancyMapPath, validMapFile)
	if err != nil {
		t.Error(err)
		return
	}

	// Test loading a valid tenancy map
	tenancies, err := LoadTenancyMap()
	assert.NoError(t, err)
	assert.Len(t, tenancies, 1)
	assert.Equal(t, "test-tenancy", tenancies[0].Tenancy)
	assert.Equal(t, "ocid1.tenancy.oc1..test", tenancies[0].TenancyID)

	// Test loading an invalid tenancy map (invalid YAML)
	invalidMapFile := filepath.Join(tempDir, "invalid-tenancy-map.yaml")
	err = os.WriteFile(invalidMapFile, []byte("invalid yaml: ]["), 0644)
	require.NoError(t, err)
	err = os.Setenv(flags.EnvKeyTenancyMapPath, invalidMapFile)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = LoadTenancyMap()
	assert.Error(t, err)

	// Test with a non-existent file
	nonExistentFile := filepath.Join(tempDir, "non-existent.yaml")
	err = os.Setenv(flags.EnvKeyTenancyMapPath, nonExistentFile)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = LoadTenancyMap()
	assert.Error(t, err)
}

// TestLookupTenancyID tests the LookupTenancyID function
func TestLookupTenancyID(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "oci-config-test")
	require.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Error(err)
		}
	}(tempDir)

	// Create a valid tenancy map file
	validMapContent := `
- environment: test
  tenancy: test-tenancy
  tenancy_id: ocid1.tenancy.oc1..test
  realm: test-realm
  compartments: 
    - test-compartments 
    - test-compartment-2
  regions: 
    - test-regions 
    - test-region-2
`
	validMapFile := filepath.Join(tempDir, "valid-tenancy-map.yaml")
	err = os.WriteFile(validMapFile, []byte(validMapContent), 0644)
	require.NoError(t, err)

	// Set the environment variable to point to the test file
	err = os.Setenv(flags.EnvKeyTenancyMapPath, validMapFile)
	if err != nil {
		t.Error(err)
		return
	}

	// Test looking up an existing tenancy
	tenancyID, err := LookupTenancyID("test-tenancy")
	assert.NoError(t, err)
	assert.Equal(t, "ocid1.tenancy.oc1..test", tenancyID)

	// Test looking up a non-existent tenancy
	_, err = LookupTenancyID("non-existent-tenancy")
	assert.Error(t, err)
}

// TestLoadOCIConfig tests the LoadOCIConfig function
func TestLoadOCIConfig(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Test with the default profile
	err := os.Unsetenv(flags.EnvKeyProfile)
	if err != nil {
		t.Error(err)
		return
	}
	provider := LoadOCIConfig()
	assert.IsType(t, common.DefaultConfigProvider(), provider)

	// Test with custom profile
	// Note: This test is limited because we can't easily verify the provider's behavior
	// without actually reading the OCI config file
	err = os.Setenv(flags.EnvKeyProfile, "CUSTOM_PROFILE")
	if err != nil {
		t.Error(err)
		return
	}
	provider = LoadOCIConfig()
	assert.NotNil(t, provider)
}

// TestUserHomeDir tests the getUserHomeDir function
func TestUserHomeDir(t *testing.T) {
	// This is a simple test to ensure the function doesn't panic
	// We can't easily test the actual return value as it depends on the system
	dir, err := getUserHomeDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)
}

// TestGetTenancyOCID tests the GetTenancyOCID function
func TestGetTenancyOCID(t *testing.T) {
	// Save the original mock function
	originalMock := MockGetTenancyOCID
	defer func() {
		// Restore original mock function
		MockGetTenancyOCID = originalMock
	}()

	// Set up a mock function for testing
	expectedTenancyID := "mock-tenancy-ocid-for-test"
	MockGetTenancyOCID = func() (string, error) {
		return expectedTenancyID, nil
	}

	// Test successful retrieval
	tenancyID, err := GetTenancyOCID()
	assert.NoError(t, err)
	assert.Equal(t, expectedTenancyID, tenancyID)

	// Test error case
	expectedError := fmt.Errorf("mock error")
	MockGetTenancyOCID = func() (string, error) {
		return "", expectedError
	}

	_, err = GetTenancyOCID()
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
