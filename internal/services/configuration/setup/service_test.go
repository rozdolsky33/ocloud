package setup

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

// TestValidateRealm tests the validateRealm function
func TestValidateRealm(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid realm OC1",
			input:       "OC1",
			expected:    "OC1",
			expectError: false,
		},
		{
			name:        "Valid realm oc1 (lowercase)",
			input:       "oc1",
			expected:    "OC1",
			expectError: false,
		},
		{
			name:        "Valid realm OC2",
			input:       "OC2",
			expected:    "OC2",
			expectError: false,
		},
		{
			name:        "Invalid realm - too long",
			input:       "OC123",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid realm - doesn't start with OC",
			input:       "AB1",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid realm - empty",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validateRealm(tc.input)
			if tc.expectError {
				assert.Error(t, err, "Expected an error for input: %s", tc.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tc.input)
				assert.Equal(t, tc.expected, result, "Expected %s but got %s", tc.expected, result)
			}
		})
	}
}

// TestValidateTenancyID tests the validateTenancyID function
func TestValidateTenancyID(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid tenancy ID with 'tenancy' in it",
			input:       "ocid1.tenancy.oc1..aaaaaaaabcdefg",
			expected:    "ocid1.tenancy.oc1..aaaaaaaabcdefg",
			expectError: false,
		},
		{
			name:        "Invalid tenancy ID - doesn't contain 'tenancy'",
			input:       "ocid1.user.oc1..aaaaaaaabcdefg",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid tenancy ID - empty",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validateTenancyID(tc.input)
			if tc.expectError {
				assert.Error(t, err, "Expected an error for input: %s", tc.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tc.input)
				assert.Equal(t, tc.expected, result, "Expected %s but got %s", tc.expected, result)
			}
		})
	}
}

// TestPromptWithValidation tests the promptWithValidation function
// This is a limited test since we can't easily test interactive functions
func TestPromptWithValidation(t *testing.T) {
	// Skip this test in normal test runs since it's interactive
	t.Skip("Skipping interactive test")

	// In a real test environment, we would mock the reader and test the function
	// For example:
	// mockReader := &MockReader{
	//     ReadStringFunc: func(delim byte) (string, error) {
	//         return "OC1\n", nil
	//     },
	// }
	// result := promptWithValidation(mockReader, "Realm", validateRealm)
	// assert.Equal(t, "OC1", result)
}