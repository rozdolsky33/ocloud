package display

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestJWTToken creates a JWT token with the given expiration time for testing.
// JWT format: header.payload.signature (we only need the payload for our tests)
func createTestJWTToken(exp int64) string {
	header := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	claims := jwtClaims{Exp: exp}
	claimsJSON, _ := json.Marshal(claims)
	payload := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(claimsJSON)
	signature := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte("fake-signature"))
	return header + "." + payload + "." + signature
}

func TestSessionExpiryFromTokenFile(t *testing.T) {
	tests := []struct {
		name        string
		tokenData   string
		expectError bool
		errorMsg    string
		checkExp    func(t *testing.T, exp time.Time)
	}{
		{
			name:        "valid token with future expiry",
			tokenData:   createTestJWTToken(time.Now().Add(1 * time.Hour).Unix()),
			expectError: false,
			checkExp: func(t *testing.T, exp time.Time) {
				assert.True(t, exp.After(time.Now()), "expiry should be in the future")
			},
		},
		{
			name:        "valid token with past expiry",
			tokenData:   createTestJWTToken(time.Now().Add(-1 * time.Hour).Unix()),
			expectError: false,
			checkExp: func(t *testing.T, exp time.Time) {
				assert.True(t, exp.Before(time.Now()), "expiry should be in the past")
			},
		},
		{
			name:        "valid token with specific timestamp",
			tokenData:   createTestJWTToken(1702400000), // 2023-12-12 18:13:20 UTC
			expectError: false,
			checkExp: func(t *testing.T, exp time.Time) {
				assert.Equal(t, int64(1702400000), exp.Unix())
			},
		},
		{
			name:        "invalid token format - no dots",
			tokenData:   "invalidtoken",
			expectError: true,
			errorMsg:    "invalid token format",
		},
		{
			name:        "invalid token format - only one part",
			tokenData:   "header",
			expectError: true,
			errorMsg:    "invalid token format",
		},
		{
			name:        "invalid base64 payload",
			tokenData:   "header.!!!invalid-base64!!!.signature",
			expectError: true,
			errorMsg:    "decode payload",
		},
		{
			name:        "invalid JSON payload",
			tokenData:   "header." + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte("not-json")) + ".sig",
			expectError: true,
			errorMsg:    "unmarshal payload",
		},
		{
			name:        "missing exp claim",
			tokenData:   "header." + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(`{"sub":"test"}`)) + ".sig",
			expectError: true,
			errorMsg:    "no exp claim",
		},
		{
			name:        "empty token file",
			tokenData:   "",
			expectError: true,
			errorMsg:    "invalid token format",
		},
		{
			name:        "token with whitespace",
			tokenData:   "  " + createTestJWTToken(time.Now().Add(1*time.Hour).Unix()) + "  \n",
			expectError: false,
			checkExp: func(t *testing.T, exp time.Time) {
				assert.True(t, exp.After(time.Now()), "expiry should be in the future")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp file with token data
			tmpDir := t.TempDir()
			tokenPath := filepath.Join(tmpDir, "token")
			err := os.WriteFile(tokenPath, []byte(tt.tokenData), 0600)
			require.NoError(t, err)

			// Test the function
			exp, err := sessionExpiryFromTokenFile(tokenPath)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkExp != nil {
					tt.checkExp(t, exp)
				}
			}
		})
	}
}

func TestSessionExpiryFromTokenFile_FileNotFound(t *testing.T) {
	_, err := sessionExpiryFromTokenFile("/nonexistent/path/token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read token file")
}

func TestGetSecurityTokenFile(t *testing.T) {
	tests := []struct {
		name           string
		profile        string
		configContent  string
		expectedPath   string
		expectFallback bool
	}{
		{
			name:    "profile with security_token_file",
			profile: "TEST",
			configContent: `[DEFAULT]
tenancy=ocid1.tenancy.default

[TEST]
fingerprint=aa:bb:cc
key_file=/path/to/key.pem
tenancy=ocid1.tenancy.test
region=us-ashburn-1
security_token_file=/custom/path/to/token
`,
			expectedPath:   "/custom/path/to/token",
			expectFallback: false,
		},
		{
			name:    "profile without security_token_file falls back",
			profile: "NOTOKEN",
			configContent: `[NOTOKEN]
fingerprint=aa:bb:cc
key_file=/path/to/key.pem
tenancy=ocid1.tenancy.test
region=us-ashburn-1
`,
			expectFallback: true,
		},
		{
			name:    "profile not found falls back",
			profile: "NONEXISTENT",
			configContent: `[DEFAULT]
tenancy=ocid1.tenancy.default
`,
			expectFallback: true,
		},
		{
			name:    "security_token_file with spaces around equals",
			profile: "SPACES",
			configContent: `[SPACES]
security_token_file = /path/with/spaces/token
`,
			expectedPath:   "/path/with/spaces/token",
			expectFallback: false,
		},
		{
			name:    "multiple profiles - finds correct one",
			profile: "MIDDLE",
			configContent: `[FIRST]
security_token_file=/first/token

[MIDDLE]
security_token_file=/middle/token

[LAST]
security_token_file=/last/token
`,
			expectedPath:   "/middle/token",
			expectFallback: false,
		},
		{
			name:    "DEFAULT profile",
			profile: "DEFAULT",
			configContent: `[DEFAULT]
fingerprint=aa:bb:cc
security_token_file=/default/session/token
`,
			expectedPath:   "/default/session/token",
			expectFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp home directory structure
			tmpHome := t.TempDir()
			ociDir := filepath.Join(tmpHome, flags.OCIConfigDirName)
			err := os.MkdirAll(ociDir, 0755)
			require.NoError(t, err)

			// Write config file
			configPath := filepath.Join(ociDir, flags.OCIConfigFileName)
			err = os.WriteFile(configPath, []byte(tt.configContent), 0600)
			require.NoError(t, err)

			// Override home directory for test
			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", tmpHome)
			defer os.Setenv("HOME", originalHome)

			// Test the function
			path, err := getSecurityTokenFile(tt.profile)
			require.NoError(t, err)

			if tt.expectFallback {
				expectedFallback := filepath.Join(tmpHome, flags.OCIConfigDirName, flags.OCISessionsDirName, tt.profile, "token")
				assert.Equal(t, expectedFallback, path)
			} else {
				assert.Equal(t, tt.expectedPath, path)
			}
		})
	}
}

func TestGetSecurityTokenFile_ConfigNotFound(t *testing.T) {
	// Create temp home without a config file
	tmpHome := t.TempDir()

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	_, err := getSecurityTokenFile("DEFAULT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open config file")
}

func TestCheckOCISessionValidity(t *testing.T) {
	tests := []struct {
		name           string
		profile        string
		configContent  string
		tokenContent   string
		expectedSubstr string
	}{
		{
			name:    "valid session",
			profile: "VALID",
			configContent: `[VALID]
security_token_file=%s
`,
			tokenContent:   createTestJWTToken(time.Now().Add(1 * time.Hour).Unix()),
			expectedSubstr: "Valid until",
		},
		{
			name:    "expired session",
			profile: "EXPIRED",
			configContent: `[EXPIRED]
security_token_file=%s
`,
			tokenContent:   createTestJWTToken(time.Now().Add(-1 * time.Hour).Unix()),
			expectedSubstr: "Session Expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp home directory structure
			tmpHome := t.TempDir()
			ociDir := filepath.Join(tmpHome, flags.OCIConfigDirName)
			sessionsDir := filepath.Join(ociDir, flags.OCISessionsDirName, tt.profile)
			err := os.MkdirAll(sessionsDir, 0755)
			require.NoError(t, err)

			// Write a token file
			tokenPath := filepath.Join(sessionsDir, "token")
			err = os.WriteFile(tokenPath, []byte(tt.tokenContent), 0600)
			require.NoError(t, err)

			// Write a config file with token path
			configPath := filepath.Join(ociDir, flags.OCIConfigFileName)
			configContent := tt.configContent
			if tt.configContent != "" {
				configContent = filepath.Join(sessionsDir, "token")
				err = os.WriteFile(configPath, []byte("["+tt.profile+"]\nsecurity_token_file="+configContent+"\n"), 0600)
			}
			require.NoError(t, err)

			// Override home directory for test
			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", tmpHome)
			defer os.Setenv("HOME", originalHome)

			// Test the function
			result := CheckOCISessionValidity(tt.profile)
			assert.Contains(t, result, tt.expectedSubstr)
		})
	}
}

func TestCheckOCISessionValidity_TokenNotFound(t *testing.T) {
	// Create a temp home directory with config but no token
	tmpHome := t.TempDir()
	ociDir := filepath.Join(tmpHome, flags.OCIConfigDirName)
	err := os.MkdirAll(ociDir, 0755)
	require.NoError(t, err)

	// Write config pointing to non-existent token
	configPath := filepath.Join(ociDir, flags.OCIConfigFileName)
	err = os.WriteFile(configPath, []byte("[NOTOKEN]\nsecurity_token_file=/nonexistent/token\n"), 0600)
	require.NoError(t, err)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	result := CheckOCISessionValidity("NOTOKEN")
	assert.Contains(t, result, "Cannot parse session token")
}

func TestCheckOCISessionValidity_InvalidToken(t *testing.T) {
	// Create temp home directory structure
	tmpHome := t.TempDir()
	ociDir := filepath.Join(tmpHome, flags.OCIConfigDirName)
	sessionsDir := filepath.Join(ociDir, flags.OCISessionsDirName, "INVALID")
	err := os.MkdirAll(sessionsDir, 0755)
	require.NoError(t, err)

	// Write an invalid token file
	tokenPath := filepath.Join(sessionsDir, "token")
	err = os.WriteFile(tokenPath, []byte("invalid-token-content"), 0600)
	require.NoError(t, err)

	// Write a config file
	configPath := filepath.Join(ociDir, flags.OCIConfigFileName)
	err = os.WriteFile(configPath, []byte("[INVALID]\nsecurity_token_file="+tokenPath+"\n"), 0600)
	require.NoError(t, err)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	result := CheckOCISessionValidity("INVALID")
	assert.Contains(t, result, "Cannot parse session token")
}

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

func TestJWTClaims(t *testing.T) {
	// Test that jwtClaims struct correctly unmarshals JSON
	tests := []struct {
		name     string
		json     string
		expected int64
	}{
		{
			name:     "basic exp claim",
			json:     `{"exp":1702400000}`,
			expected: 1702400000,
		},
		{
			name:     "exp with other claims",
			json:     `{"exp":1702400000,"sub":"user123","iss":"oracle"}`,
			expected: 1702400000,
		},
		{
			name:     "zero exp",
			json:     `{"exp":0}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var claims jwtClaims
			err := json.Unmarshal([]byte(tt.json), &claims)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, claims.Exp)
		})
	}
}
