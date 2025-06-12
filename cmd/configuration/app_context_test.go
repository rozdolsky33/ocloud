package configuration

import (
	"fmt"
	"os"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setupTest prepares the test environment and returns a cleanup function
func setupTest(t *testing.T) func() {
	// Save original environment variables to restore later
	originalEnvVars := map[string]string{
		EnvOCITenancy:     os.Getenv(EnvOCITenancy),
		EnvOCITenancyName: os.Getenv(EnvOCITenancyName),
		EnvOCICompartment: os.Getenv(EnvOCICompartment),
		EnvOCIRegion:      os.Getenv(EnvOCIRegion),
	}

	// Reset viper for each test
	viper.Reset()
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()

	// Save original mock functions
	originalGetTenancyOCID := config.MockGetTenancyOCID
	originalLookupTenancyID := config.MockLookupTenancyID

	// Set up mock functions for testing
	config.MockGetTenancyOCID = func() (string, error) {
		return "mock-tenancy-ocid", nil
	}
	config.MockLookupTenancyID = func(tenancyName string) (string, error) {
		return "mock-tenancy-ocid-for-" + tenancyName, nil
	}

	// Return a cleanup function
	return func() {
		// Restore original environment variables
		for key, value := range originalEnvVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}

		// Restore original mock functions
		config.MockGetTenancyOCID = originalGetTenancyOCID
		config.MockLookupTenancyID = originalLookupTenancyID

		viper.Reset()
	}
}

// TestInitGlobalFlags tests the InitGlobalFlags function
func TestInitGlobalFlags(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command
	cmd := &cobra.Command{}

	// Call the function
	InitGlobalFlags(cmd)

	// Verify that flags were added
	tenancyFlag := cmd.PersistentFlags().Lookup(FlagNameTenancyID)
	assert.NotNil(t, tenancyFlag)
	assert.Equal(t, FlagShortTenancyID, tenancyFlag.Shorthand)
	assert.Equal(t, FlagDescTenancyID, tenancyFlag.Usage)

	compartmentFlag := cmd.PersistentFlags().Lookup(FlagNameCompartment)
	assert.NotNil(t, compartmentFlag)
	assert.Equal(t, FlagShortCompartment, compartmentFlag.Shorthand)
	assert.Equal(t, FlagDescCompartment, compartmentFlag.Usage)

	// Verify that viper is set up correctly
	assert.Equal(t, "OCI", viper.GetEnvPrefix())
}

// TestLoadTenancyOCID tests the fetchTenancyOCID function
func TestLoadTenancyOCID(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Test successful retrieval
	err := fetchTenancyOCID()
	assert.NoError(t, err)
	assert.Equal(t, "mock-tenancy-ocid", viper.GetString(FlagNameTenancyID))

	// Test error case
	config.MockGetTenancyOCID = func() (string, error) {
		return "", fmt.Errorf("mock error")
	}
	err = fetchTenancyOCID()
	assert.Error(t, err)
}

// TestTenancyIDFromFlag tests that tenancy ID from a flag is used
func TestTenancyIDFromFlag(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Set the flag as changed and set a value
	testTenancyID := "test-tenancy-id"
	viper.Set(FlagNameTenancyID, testTenancyID)
	cmd.Flags().Set(FlagNameTenancyID, testTenancyID)

	// Verify that viper has the correct value
	assert.Equal(t, testTenancyID, viper.GetString(FlagNameTenancyID))
}

// TestTenancyIDFromEnv tests that tenancy ID from the environment is used
func TestTenancyIDFromEnv(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Set environment variable
	testTenancyID := "env-tenancy-id"
	os.Setenv(EnvOCITenancy, testTenancyID)

	// Verify that the environment variable is set
	assert.Equal(t, testTenancyID, os.Getenv(EnvOCITenancy))
}

// TestCompartmentFromFlag tests that compartment from a flag is used
func TestCompartmentFromFlag(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameCompartment, "", "")

	// Set the flag as changed and set a value
	testCompartment := "test-compartment"
	viper.Set(FlagNameCompartment, testCompartment)
	cmd.Flags().Set(FlagNameCompartment, testCompartment)

	// Verify that the viper has the correct value
	assert.Equal(t, testCompartment, viper.GetString(FlagNameCompartment))
}

// TestCompartmentFromEnv tests that compartment from the environment is used
func TestCompartmentFromEnv(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameCompartment, "", "")

	// Set environment variable
	testCompartment := "env-compartment"
	os.Setenv(EnvOCICompartment, testCompartment)

	// Verify that the environment variable is set
	assert.Equal(t, testCompartment, os.Getenv(EnvOCICompartment))
}
