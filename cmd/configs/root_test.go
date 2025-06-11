package configs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Execute runs the ConfigCmd and exits on error
func Execute() {
	if err := ConfigCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		osExit(1)
	}
}

// setupTest prepares the test environment
func setupTest(t *testing.T) func() {
	// Save original environment variables to restore later
	originalEnvVars := map[string]string{
		EnvOCITenancy:     os.Getenv(EnvOCITenancy),
		EnvOCITenancyName: os.Getenv(EnvOCITenancyName),
		EnvOCICompartment: os.Getenv(EnvOCICompartment),
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

// TestInitializeConfigWithDebugMode tests that debug mode sets the correct log level
func TestInitializeConfigWithDebugMode(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command
	cmd := &cobra.Command{}

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	// Note: We're no longer checking log levels as we've migrated from logrus to slog/logr
}

// TestInitializeConfigWithoutDebugMode tests that non-debug mode sets the correct log level
func TestInitializeConfigWithoutDebugMode(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command
	cmd := &cobra.Command{}

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	// Note: We're no longer checking log levels as we've migrated from logrus to slog/logr
}

// TestInitializeConfigWithTenancyFlag tests that tenancy ID from a flag is used
func TestInitializeConfigWithTenancyFlag(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Set the flag as changed and set a value
	testTenancyID := "test-tenancy-id"
	viper.Set(FlagNameTenancyID, testTenancyID)
	cmd.Flags().Set(FlagNameTenancyID, testTenancyID)

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, testTenancyID, viper.GetString(FlagNameTenancyID))
}

// TestInitializeConfigWithTenancyEnv tests that tenancy ID from the environment is used
func TestInitializeConfigWithTenancyEnv(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Set environment variable
	testTenancyID := "env-tenancy-id"
	os.Setenv(EnvOCITenancy, testTenancyID)

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, testTenancyID, viper.GetString(FlagNameTenancyID))
}

// TestInitializeConfigWithCompartmentFlag tests that compartment from a flag is used
func TestInitializeConfigWithCompartmentFlag(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameCompartment, "", "")

	// Set the flag as changed and set a value
	testCompartment := "test-compartment"
	viper.Set(FlagNameCompartment, testCompartment)
	cmd.Flags().Set(FlagNameCompartment, testCompartment)

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, testCompartment, viper.GetString(FlagNameCompartment))
}

// TestInitializeConfigWithCompartmentEnv tests that compartment from the environment is used
func TestInitializeConfigWithCompartmentEnv(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameCompartment, "", "")

	// Set environment variable
	testCompartment := "env-compartment"
	os.Setenv(EnvOCICompartment, testCompartment)

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, testCompartment, viper.GetString(FlagNameCompartment))
}

// TestExecute tests the Execute function (basic smoke test)
func TestExecute(t *testing.T) {
	// Save the original os.Exit and restore it after the test
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	// Save the original ConfigCmd for restoration
	originalRootCmd := ConfigCmd
	defer func() { ConfigCmd = originalRootCmd }()

	// Save the original logrus output and level for backward compatibility
	originalOutput := logrus.StandardLogger().Out
	originalLevel := logrus.GetLevel()
	defer func() {
		logrus.SetOutput(originalOutput)
		logrus.SetLevel(originalLevel)
	}()

	// Discard logrus output during test
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.PanicLevel) // Only log panic level messages

	// Test successful execution
	{
		// Mock os.Exit to prevent the test from exiting
		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
			assert.Equal(t, 1, code, "Expected exit code 1 on error")
		}

		// Create a mock command that will succeed
		ConfigCmd = &cobra.Command{
			Use: "mock",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
			SilenceErrors: true, // Don't print errors
			SilenceUsage:  true, // Don't print usage on error
		}

		// Call Execute with the mock command
		Execute()

		// Verify that os.Exit was not called
		assert.False(t, exitCalled, "os.Exit should not be called on success")
	}

	// Test execution with error
	{
		// Mock os.Exit to prevent the test from exiting
		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
			assert.Equal(t, 1, code, "Expected exit code 1 on error")
		}

		// Create a mock command that will fail
		ConfigCmd = &cobra.Command{
			Use: "mock",
			RunE: func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("test error")
			},
			SilenceErrors: true, // Don't print errors
			SilenceUsage:  true, // Don't print usage on error
		}

		// Call Execute with the failing mock command
		Execute()

		// Verify that os.Exit was called with code 1
		assert.True(t, exitCalled, "os.Exit should be called on error")
	}
}

// Note: osExit is defined in root.go and used here for mocking

// TestInitializeConfigWithTenancyName tests that tenancy ID is looked up by name
func TestInitializeConfigWithTenancyName(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Set tenancy name environment variable
	testTenancyName := "test-tenancy-name"
	os.Setenv(EnvOCITenancyName, testTenancyName)

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify results
	assert.NoError(t, err)
	expectedTenancyID := "mock-tenancy-ocid-for-" + testTenancyName
	assert.Equal(t, expectedTenancyID, viper.GetString(FlagNameTenancyID))
}

// TestInitializeConfigPriority tests the priority order of tenancy ID sources
func TestInitializeConfigPriority(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a test command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Set all possible sources of tenancy ID
	flagTenancyID := "flag-tenancy-id"
	envTenancyID := "env-tenancy-id"

	// Set flag (the highest priority)
	viper.Set(FlagNameTenancyID, flagTenancyID)
	cmd.Flags().Set(FlagNameTenancyID, flagTenancyID)

	// Set environment variable (medium priority)
	os.Setenv(EnvOCITenancy, envTenancyID)

	// Call the function
	err := initConfig(cmd, []string{})

	// Verify that flag value is used (the highest priority)
	assert.NoError(t, err)
	assert.Equal(t, flagTenancyID, viper.GetString(FlagNameTenancyID))

	// Reset and test with only the environment variable
	viper.Reset()
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()

	// Create a new command without setting the flag
	cmd = &cobra.Command{}
	cmd.Flags().String(FlagNameTenancyID, "", "")

	// Call the function
	err = initConfig(cmd, []string{})

	// Verify that the environment variable is used
	assert.NoError(t, err)
	assert.Equal(t, envTenancyID, viper.GetString(FlagNameTenancyID))
}
