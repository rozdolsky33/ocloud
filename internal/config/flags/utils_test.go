package flags

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetBoolFlag(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    bool
		flagExists   bool
		defaultValue bool
		expected     bool
	}{
		{
			name:         "flag exists and is true",
			flagValue:    true,
			flagExists:   true,
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "flag exists and is false",
			flagValue:    false,
			flagExists:   true,
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "flag does not exist - return default true",
			flagExists:   false,
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "flag does not exist - return default false",
			flagExists:   false,
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flagName := "test-flag"

			if tt.flagExists {
				cmd.Flags().Bool(flagName, tt.flagValue, "test flag")
			}

			result := GetBoolFlag(cmd, flagName, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStringFlag(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    string
		flagExists   bool
		defaultValue string
		expected     string
	}{
		{
			name:         "flag exists with value",
			flagValue:    "test-value",
			flagExists:   true,
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "flag exists with empty value",
			flagValue:    "",
			flagExists:   true,
			defaultValue: "default",
			expected:     "",
		},
		{
			name:         "flag does not exist",
			flagExists:   false,
			defaultValue: "default-value",
			expected:     "default-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flagName := "test-flag"

			if tt.flagExists {
				cmd.Flags().String(flagName, tt.flagValue, "test flag")
			}

			result := GetStringFlag(cmd, flagName, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetIntFlag(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    int
		flagExists   bool
		defaultValue int
		expected     int
	}{
		{
			name:         "flag exists with positive value",
			flagValue:    42,
			flagExists:   true,
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "flag exists with zero value",
			flagValue:    0,
			flagExists:   true,
			defaultValue: 10,
			expected:     0,
		},
		{
			name:         "flag exists with negative value",
			flagValue:    -5,
			flagExists:   true,
			defaultValue: 10,
			expected:     -5,
		},
		{
			name:         "flag does not exist",
			flagExists:   false,
			defaultValue: 100,
			expected:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flagName := "test-flag"

			if tt.flagExists {
				cmd.Flags().Int(flagName, tt.flagValue, "test flag")
			}

			result := GetIntFlag(cmd, flagName, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStringSliceFlag(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    []string
		flagExists   bool
		defaultValue []string
		expected     []string
	}{
		{
			name:         "flag exists with multiple values",
			flagValue:    []string{"value1", "value2", "value3"},
			flagExists:   true,
			defaultValue: []string{"default"},
			expected:     []string{"value1", "value2", "value3"},
		},
		{
			name:         "flag exists with single value",
			flagValue:    []string{"single"},
			flagExists:   true,
			defaultValue: []string{"default1", "default2"},
			expected:     []string{"single"},
		},
		{
			name:         "flag exists with empty slice",
			flagValue:    []string{},
			flagExists:   true,
			defaultValue: []string{"default"},
			expected:     []string{},
		},
		{
			name:         "flag does not exist",
			flagExists:   false,
			defaultValue: []string{"default1", "default2"},
			expected:     []string{"default1", "default2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flagName := "test-flag"

			if tt.flagExists {
				cmd.Flags().StringSlice(flagName, tt.flagValue, "test flag")
			}

			result := GetStringSliceFlag(cmd, flagName, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFloat64Flag(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    float64
		flagExists   bool
		defaultValue float64
		expected     float64
	}{
		{
			name:         "flag exists with positive value",
			flagValue:    3.14159,
			flagExists:   true,
			defaultValue: 0.0,
			expected:     3.14159,
		},
		{
			name:         "flag exists with zero value",
			flagValue:    0.0,
			flagExists:   true,
			defaultValue: 1.5,
			expected:     0.0,
		},
		{
			name:         "flag exists with negative value",
			flagValue:    -2.71828,
			flagExists:   true,
			defaultValue: 1.0,
			expected:     -2.71828,
		},
		{
			name:         "flag does not exist",
			flagExists:   false,
			defaultValue: 9.99,
			expected:     9.99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flagName := "test-flag"

			if tt.flagExists {
				cmd.Flags().Float64(flagName, tt.flagValue, "test flag")
			}

			result := GetFloat64Flag(cmd, flagName, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDurationFlag(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    int64
		flagExists   bool
		defaultValue int64
		expected     int64
	}{
		{
			name:         "flag exists with positive duration",
			flagValue:    int64(time.Hour),
			flagExists:   true,
			defaultValue: int64(time.Minute),
			expected:     int64(time.Hour),
		},
		{
			name:         "flag exists with zero duration",
			flagValue:    0,
			flagExists:   true,
			defaultValue: int64(time.Second),
			expected:     0,
		},
		{
			name:         "flag does not exist",
			flagExists:   false,
			defaultValue: int64(time.Minute * 30),
			expected:     int64(time.Minute * 30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flagName := "test-flag"

			if tt.flagExists {
				cmd.Flags().Int64(flagName, tt.flagValue, "test flag")
			}

			result := GetDurationFlag(cmd, flagName, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlagErrorHandling(t *testing.T) {
	t.Run("GetBoolFlag handles flag type mismatch", func(t *testing.T) {
		cmd := &cobra.Command{}
		// Add a string flag but try to read it as bool
		cmd.Flags().String("test-flag", "string-value", "test")

		result := GetBoolFlag(cmd, "test-flag", true)
		assert.True(t, result) // Should return default value on error
	})

	t.Run("GetStringFlag handles flag type mismatch", func(t *testing.T) {
		cmd := &cobra.Command{}
		// Add a bool flag but try to read it as string
		cmd.Flags().Bool("test-flag", true, "test")

		result := GetStringFlag(cmd, "test-flag", "default")
		assert.Equal(t, "default", result) // Should return default value on error
	})

	t.Run("GetIntFlag handles flag type mismatch", func(t *testing.T) {
		cmd := &cobra.Command{}
		// Add a string flag but try to read it as int
		cmd.Flags().String("test-flag", "not-a-number", "test")

		result := GetIntFlag(cmd, "test-flag", 42)
		assert.Equal(t, 42, result) // Should return default value on error
	})
}
