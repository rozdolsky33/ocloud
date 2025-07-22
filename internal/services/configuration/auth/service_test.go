package auth

import (
	"bytes"
	"crypto/rsa"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	// Create a mock logger
	mockLogger := logr.Discard()

	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger: mockLogger,
	}

	// Create a new service
	service := NewService(appCtx)

	// Verify that the service was created correctly
	assert.NotNil(t, service)
	assert.Equal(t, appCtx, service.appCtx)
	assert.Equal(t, mockLogger, service.logger)
}

func TestGetOCIRegions(t *testing.T) {
	// Create a mock logger
	mockLogger := logr.Discard()

	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger: mockLogger,
	}

	// Create a new service
	service := NewService(appCtx)

	// Get the regions
	regions := service.GetOCIRegions()

	// Verify that the regions were returned correctly
	assert.NotEmpty(t, regions)
	assert.Equal(t, "1", regions[0].ID)
	assert.Equal(t, "af-johannesburg-1", regions[0].Name)

	// Verify that all regions have an ID and a name
	for _, region := range regions {
		assert.NotEmpty(t, region.ID)
		assert.NotEmpty(t, region.Name)
	}
}

func TestGetCurrentEnvironment(t *testing.T) {
	// Save original environment variables to restore later
	originalProfile := os.Getenv("OCI_CLI_PROFILE")
	originalRegion := os.Getenv("OCI_REGION")
	defer func() {
		os.Setenv("OCI_CLI_PROFILE", originalProfile)
		os.Setenv("OCI_REGION", originalRegion)
	}()

	// Set test environment variables
	os.Setenv("OCI_CLI_PROFILE", "TEST_PROFILE")
	os.Setenv("OCI_REGION", "us-ashburn-1")

	// Create a mock provider
	mockProvider := oci.NewMockConfigurationProvider()

	// Save original mock functions to restore later
	originalLoadTenancyMap := config.MockLoadTenancyMap
	defer func() { config.MockLoadTenancyMap = originalLoadTenancyMap }()

	// Set up mock for LoadTenancyMap
	config.MockLoadTenancyMap = func() ([]config.MappingsFile, error) {
		return []config.MappingsFile{
			{
				Tenancy:   "example-tenancy",
				TenancyID: "ocid1.tenancy.oc1..mock-tenancy-id",
			},
		}, nil
	}

	// Create a mock logger
	mockLogger := logger.NewTestLogger()

	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger:   mockLogger,
		Provider: mockProvider,
	}

	// Create a new service
	service := NewService(appCtx)

	// Get the current environment
	result, err := service.GetCurrentEnvironment()

	// Verify that the result was returned correctly
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ocid1.tenancy.oc1..mock-tenancy-id", result.TenancyID)
	assert.Equal(t, "example-tenancy", result.TenancyName)
	assert.Equal(t, "TEST_PROFILE", result.Profile)
	assert.Equal(t, "us-ashburn-1", result.Region)
}

// ErrorMockProvider is a mock provider that returns an error for TenancyOCID
type ErrorMockProvider struct {
	common.ConfigurationProvider
}

func (p *ErrorMockProvider) TenancyOCID() (string, error) {
	return "", assert.AnError
}

func (p *ErrorMockProvider) UserOCID() (string, error) {
	return "ocid1.user.oc1..mock-user-id", nil
}

func (p *ErrorMockProvider) KeyFingerprint() (string, error) {
	return "mock-key-fingerprint", nil
}

func (p *ErrorMockProvider) Region() (string, error) {
	return "us-ashburn-1", nil
}

func (p *ErrorMockProvider) KeyID() (string, error) {
	return "mock-key-id", nil
}

func (p *ErrorMockProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	return nil, nil
}

func (p *ErrorMockProvider) Passphrase() (string, error) {
	return "", nil
}

func (p *ErrorMockProvider) AuthType() (common.AuthConfig, error) {
	return common.AuthConfig{
		AuthType:         common.UserPrincipal,
		IsFromConfigFile: false,
	}, nil
}

func TestGetCurrentEnvironment_Error(t *testing.T) {
	// Create a mock logger
	mockLogger := logger.NewTestLogger()

	// Create a mock application context with a provider that returns an error for TenancyOCID
	appCtx := &app.ApplicationContext{
		Logger:   mockLogger,
		Provider: &ErrorMockProvider{},
	}

	// Create a new service
	service := NewService(appCtx)

	// Get the current environment
	result, err := service.GetCurrentEnvironment()

	// Verify that an error was returned
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPromptYesNo(t *testing.T) {
	// Save stdin to restore later
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Yes",
			input:    "y\n",
			expected: true,
		},
		{
			name:     "Yes uppercase",
			input:    "Y\n",
			expected: true,
		},
		{
			name:     "Yes full",
			input:    "yes\n",
			expected: true,
		},
		{
			name:     "No",
			input:    "n\n",
			expected: false,
		},
		{
			name:     "No uppercase",
			input:    "N\n",
			expected: false,
		},
		{
			name:     "No full",
			input:    "no\n",
			expected: false,
		},
		{
			name:     "Invalid input then yes",
			input:    "invalid\ny\n",
			expected: true,
		},
		{
			name:     "Invalid input then no",
			input:    "invalid\nn\n",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pipe to provide input
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write the input to the pipe
			go func() {
				_, _ = w.Write([]byte(tc.input))
				w.Close()
			}()

			// Save stdout to restore later
			oldStdout := os.Stdout
			// Create a pipe to capture stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			// Call the function
			result := promptYesNo("Test question")

			// Restore stdout
			os.Stdout = oldStdout
			wOut.Close()

			// Verify the result
			assert.Equal(t, tc.expected, result)

			// Read the output to clear the buffer
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(rOut)
		})
	}
}

func TestPromptForProfile(t *testing.T) {
	// Save stdin to restore later
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a mock logger
	mockLogger := logger.NewTestLogger()

	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger: mockLogger,
	}

	// Create a new service
	service := NewService(appCtx)

	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Default profile",
			input:    "1\n",
			expected: "DEFAULT",
		},
		{
			name:     "Custom profile",
			input:    "2\ncustom-profile\n",
			expected: "custom-profile",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pipe to provide input
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write the input to the pipe
			go func() {
				_, _ = w.Write([]byte(tc.input))
				w.Close()
			}()

			// Save stdout to restore later
			oldStdout := os.Stdout
			// Create a pipe to capture stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			// Call the function
			result, err := service.PromptForProfile()

			// Restore stdout
			os.Stdout = oldStdout
			wOut.Close()

			// Verify the result
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)

			// Read the output to clear the buffer
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(rOut)
		})
	}
}

func TestPromptForRegion(t *testing.T) {
	// Save stdin to restore later
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a mock logger
	mockLogger := logger.NewTestLogger()

	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger: mockLogger,
	}

	// Create a new service
	service := NewService(appCtx)

	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Region by index",
			input:    "1\n",
			expected: "af-johannesburg-1",
		},
		{
			name:     "Region by name",
			input:    "us-ashburn-1\n",
			expected: "us-ashburn-1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pipe to provide input
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write the input to the pipe
			go func() {
				_, _ = w.Write([]byte(tc.input))
				w.Close()
			}()

			// Save stdout to restore later
			oldStdout := os.Stdout
			// Create a pipe to capture stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			// Call the function
			result, err := service.PromptForRegion()

			// Restore stdout
			os.Stdout = oldStdout
			wOut.Close()

			// Verify the result
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)

			// Read the output to clear the buffer
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(rOut)
		})
	}
}

// Note: We can't easily test the Authenticate function because it runs an external command
// In a real-world scenario, we would refactor the function to accept a command runner interface
// that could be mocked for testing. For now, we'll skip this test.
