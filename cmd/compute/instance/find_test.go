package instance

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	instaceFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// Mock application context for testing
func createMockAppContext() *app.ApplicationContext {
	return &app.ApplicationContext{
		CompartmentName: "test-compartment",
		// Add other necessary fields as needed
	}
}

// Helper function to execute a command and capture output
func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	
	err := root.Execute()
	return strings.TrimSpace(buf.String()), err
}

// TestNewFindCmd tests the creation and configuration of the find command
func TestNewFindCmd(t *testing.T) {
	tests := []struct {
		name     string
		appCtx   *app.ApplicationContext
		wantUse  string
		wantShort string
	}{
		{
			name:      "creates find command with correct configuration",
			appCtx:    createMockAppContext(),
			wantUse:   "find [pattern]",
			wantShort: "Find instances by name pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewFindCmd(tt.appCtx)
			
			if cmd == nil {
				t.Fatal("NewFindCmd() returned nil")
			}
			
			if cmd.Use \!= tt.wantUse {
				t.Errorf("NewFindCmd().Use = %v, want %v", cmd.Use, tt.wantUse)
			}
			
			if cmd.Short \!= tt.wantShort {
				t.Errorf("NewFindCmd().Short = %v, want %v", cmd.Short, tt.wantShort)
			}
			
			// Test aliases
			expectedAliases := []string{"f"}
			if len(cmd.Aliases) \!= len(expectedAliases) {
				t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
			} else {
				for i, alias := range expectedAliases {
					if cmd.Aliases[i] \!= alias {
						t.Errorf("Expected alias %s, got %s", alias, cmd.Aliases[i])
					}
				}
			}
			
			// Test that the command requires exactly 1 argument
			if cmd.Args == nil {
				t.Error("Expected Args to be set")
			}
			
			// Test that SilenceUsage and SilenceErrors are set
			if \!cmd.SilenceUsage {
				t.Error("Expected SilenceUsage to be true")
			}
			
			if \!cmd.SilenceErrors {
				t.Error("Expected SilenceErrors to be true")
			}
			
			// Test that Long description contains expected content
			if \!strings.Contains(cmd.Long, "fuzzy matching algorithm") {
				t.Error("Expected Long description to mention fuzzy matching algorithm")
			}
			
			// Test that Example contains expected content
			if \!strings.Contains(cmd.Example, "ocloud compute instance find") {
				t.Error("Expected Example to contain command usage examples")
			}
		})
	}
}

// TestNewFindCmd_HasFlags tests that the find command has the expected flags
func TestNewFindCmd_HasFlags(t *testing.T) {
	appCtx := createMockAppContext()
	cmd := NewFindCmd(appCtx)
	
	// Check that the all flag is added (from instaceFlags.AllInfoFlag.Add(cmd))
	allFlag := cmd.Flags().Lookup(flags.FlagNameAll)
	if allFlag == nil {
		t.Error("Expected 'all' flag to be present")
	}
}

// TestNewFindCmd_Arguments tests command argument validation
func TestNewFindCmd_Arguments(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid single argument",
			args:    []string{"web"},
			wantErr: false,
		},
		{
			name:        "no arguments provided",
			args:        []string{},
			wantErr:     true,
			errContains: "requires exactly 1 arg",
		},
		{
			name:        "too many arguments",
			args:        []string{"web", "server", "extra"},
			wantErr:     true,
			errContains: "requires exactly 1 arg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appCtx := createMockAppContext()
			cmd := NewFindCmd(appCtx)
			
			// Create a root command to test argument validation
			rootCmd := &cobra.Command{Use: "test"}
			rootCmd.AddCommand(cmd)
			
			_, err := executeCommand(rootCmd, append([]string{"find"}, tt.args...)...)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errContains \!= "" && \!strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err \!= nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestRunFindCommand tests the RunFindCommand function with various scenarios
func TestRunFindCommand(t *testing.T) {
	// Mock the instance.FindInstances function to avoid external dependencies
	originalFindInstances := instance.FindInstances
	defer func() {
		instance.FindInstances = originalFindInstances
	}()

	tests := []struct {
		name           string
		args           []string
		flagSetup      func(*cobra.Command)
		mockReturn     error
		expectPattern  string
		expectDetails  bool
		expectJSON     bool
		wantErr        bool
	}{
		{
			name: "basic pattern search",
			args: []string{"webserver"},
			flagSetup: func(cmd *cobra.Command) {
				// No additional flags
			},
			mockReturn:    nil,
			expectPattern: "webserver",
			expectDetails: false,
			expectJSON:    false,
			wantErr:       false,
		},
		{
			name: "search with all flag",
			args: []string{"database"},
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set(flags.FlagNameAll, "true")
			},
			mockReturn:    nil,
			expectPattern: "database",
			expectDetails: true,
			expectJSON:    false,
			wantErr:       false,
		},
		{
			name: "search with json flag",
			args: []string{"api"},
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set(flags.FlagNameJSON, "true")
			},
			mockReturn:    nil,
			expectPattern: "api",
			expectDetails: false,
			expectJSON:    true,
			wantErr:       false,
		},
		{
			name: "search with both flags",
			args: []string{"server"},
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set(flags.FlagNameAll, "true")
				cmd.Flags().Set(flags.FlagNameJSON, "true")
			},
			mockReturn:    nil,
			expectPattern: "server",
			expectDetails: true,
			expectJSON:    true,
			wantErr:       false,
		},
		{
			name: "service returns error",
			args: []string{"failing"},
			flagSetup: func(cmd *cobra.Command) {
				// No additional flags
			},
			mockReturn:    errors.New("service error"),
			expectPattern: "failing",
			expectDetails: false,
			expectJSON:    false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPattern string
			var capturedJSON bool
			var capturedDetails bool
			var capturedAppCtx *app.ApplicationContext

			// Mock the FindInstances function
			instance.FindInstances = func(appCtx *app.ApplicationContext, pattern string, useJSON bool, showDetails bool) error {
				capturedPattern = pattern
				capturedJSON = useJSON
				capturedDetails = showDetails
				capturedAppCtx = appCtx
				return tt.mockReturn
			}

			// Create command and set up flags
			appCtx := createMockAppContext()
			cmd := NewFindCmd(appCtx)
			
			// Add the flags that will be checked
			instaceFlags.AllInfoFlag.Add(cmd)
			cmd.Flags().Bool(flags.FlagNameJSON, false, "JSON output")
			
			tt.flagSetup(cmd)

			// Execute the command
			err := RunFindCommand(cmd, tt.args, appCtx)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Error("Expected an error but got none")
				}
			} else {
				if err \!= nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}

			// Verify the correct parameters were passed to the service
			if capturedPattern \!= tt.expectPattern {
				t.Errorf("Expected pattern %q, got %q", tt.expectPattern, capturedPattern)
			}

			if capturedJSON \!= tt.expectJSON {
				t.Errorf("Expected JSON flag %v, got %v", tt.expectJSON, capturedJSON)
			}

			if capturedDetails \!= tt.expectDetails {
				t.Errorf("Expected details flag %v, got %v", tt.expectDetails, capturedDetails)
			}

			// Verify the correct app context was passed
			if capturedAppCtx \!= appCtx {
				t.Error("Expected the same app context to be passed to service")
			}
		})
	}
}

// TestRunFindCommand_EmptyPattern tests edge case with empty pattern
func TestRunFindCommand_EmptyPattern(t *testing.T) {
	// Mock the instance.FindInstances function
	originalFindInstances := instance.FindInstances
	defer func() {
		instance.FindInstances = originalFindInstances
	}()

	var capturedPattern string
	instance.FindInstances = func(appCtx *app.ApplicationContext, pattern string, useJSON bool, showDetails bool) error {
		capturedPattern = pattern
		return nil
	}

	appCtx := createMockAppContext()
	cmd := NewFindCmd(appCtx)
	instaceFlags.AllInfoFlag.Add(cmd)

	// Test with empty string pattern
	err := RunFindCommand(cmd, []string{""}, appCtx)

	if err \!= nil {
		t.Errorf("Expected no error with empty pattern, got: %v", err)
	}

	if capturedPattern \!= "" {
		t.Errorf("Expected empty pattern to be passed through, got %q", capturedPattern)
	}
}

// TestRunFindCommand_SpecialCharacters tests patterns with special characters
func TestRunFindCommand_SpecialCharacters(t *testing.T) {
	// Mock the instance.FindInstances function
	originalFindInstances := instance.FindInstances
	defer func() {
		instance.FindInstances = originalFindInstances
	}()

	tests := []struct {
		name    string
		pattern string
	}{
		{"pattern with spaces", "web server"},
		{"pattern with dashes", "api-gateway"},
		{"pattern with underscores", "db_server"},
		{"pattern with dots", "app.example.com"},
		{"pattern with numbers", "server123"},
		{"pattern with special chars", "test@#$%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPattern string
			instance.FindInstances = func(appCtx *app.ApplicationContext, pattern string, useJSON bool, showDetails bool) error {
				capturedPattern = pattern
				return nil
			}

			appCtx := createMockAppContext()
			cmd := NewFindCmd(appCtx)
			instaceFlags.AllInfoFlag.Add(cmd)

			err := RunFindCommand(cmd, []string{tt.pattern}, appCtx)

			if err \!= nil {
				t.Errorf("Expected no error with pattern %q, got: %v", tt.pattern, err)
			}

			if capturedPattern \!= tt.pattern {
				t.Errorf("Expected pattern %q, got %q", tt.pattern, capturedPattern)
			}
		})
	}
}

// TestRunFindCommand_LoggingBehavior tests that logging is called correctly
func TestRunFindCommand_LoggingBehavior(t *testing.T) {
	// Mock the instance.FindInstances function
	originalFindInstances := instance.FindInstances
	defer func() {
		instance.FindInstances = originalFindInstances
	}()

	// Mock logger to capture logging calls
	var loggedLevel logger.LogLevel
	var loggedMessage string
	var loggedFields []interface{}
	
	originalLogWithLevel := logger.LogWithLevel
	defer func() {
		logger.LogWithLevel = originalLogWithLevel
	}()

	logger.LogWithLevel = func(l *logger.Logger, level logger.LogLevel, msg string, fields ...interface{}) {
		loggedLevel = level
		loggedMessage = msg
		loggedFields = fields
	}

	instance.FindInstances = func(appCtx *app.ApplicationContext, pattern string, useJSON bool, showDetails bool) error {
		return nil
	}

	appCtx := createMockAppContext()
	cmd := NewFindCmd(appCtx)
	instaceFlags.AllInfoFlag.Add(cmd)

	err := RunFindCommand(cmd, []string{"testpattern"}, appCtx)

	if err \!= nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify logging was called with correct parameters
	if loggedLevel \!= logger.Debug {
		t.Errorf("Expected debug log level, got %v", loggedLevel)
	}

	if loggedMessage \!= "Running instance find command" {
		t.Errorf("Expected log message 'Running instance find command', got %q", loggedMessage)
	}

	// Check that logged fields contain expected keys
	expectedFields := map[string]bool{
		"pattern": false,
		"in compartment": false,
		"json": false,
	}

	for i := 0; i < len(loggedFields)-1; i += 2 {
		if key, ok := loggedFields[i].(string); ok {
			if _, exists := expectedFields[key]; exists {
				expectedFields[key] = true
			}
		}
	}

	for key, found := range expectedFields {
		if \!found {
			t.Errorf("Expected log field %q to be present", key)
		}
	}
}

// TestRunFindCommand_NilAppContext tests handling of nil application context
func TestRunFindCommand_NilAppContext(t *testing.T) {
	// Mock the instance.FindInstances function
	originalFindInstances := instance.FindInstances
	defer func() {
		instance.FindInstances = originalFindInstances
	}()

	var capturedAppCtx *app.ApplicationContext
	instance.FindInstances = func(appCtx *app.ApplicationContext, pattern string, useJSON bool, showDetails bool) error {
		capturedAppCtx = appCtx
		return nil
	}

	// Test with nil app context
	cmd := NewFindCmd(nil)
	instaceFlags.AllInfoFlag.Add(cmd)

	err := RunFindCommand(cmd, []string{"test"}, nil)

	if err \!= nil {
		t.Errorf("Expected no error with nil app context, got: %v", err)
	}

	if capturedAppCtx \!= nil {
		t.Error("Expected nil app context to be passed through")
	}
}

// TestFindCommandConstants tests that the command constants are properly defined
func TestFindCommandConstants(t *testing.T) {
	// Test that findLong contains expected documentation
	expectedContent := []string{
		"Find instances in the specified compartment",
		"fuzzy matching algorithm",
		"Searchable Fields",
		"Name: Instance name",
		"InstanceName",
		"InstanceOperatingSystem",
		"TagValues",
		"partial matches are supported",
	}

	for _, content := range expectedContent {
		if \!strings.Contains(findLong, content) {
			t.Errorf("Expected findLong to contain %q", content)
		}
	}

	// Test that findExamples contains expected usage examples
	expectedExamples := []string{
		"ocloud compute instance find web",
		"ocloud compute instance find 8.10",
		"ocloud compute instance find api --all",
		"ocloud compute instance find server --json",
		"ocloud compute instance find oracle",
	}

	for _, example := range expectedExamples {
		if \!strings.Contains(findExamples, example) {
			t.Errorf("Expected findExamples to contain %q", example)
		}
	}
}

// BenchmarkNewFindCmd benchmarks the command creation
func BenchmarkNewFindCmd(b *testing.B) {
	appCtx := createMockAppContext()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewFindCmd(appCtx)
	}
}

// BenchmarkRunFindCommand benchmarks the command execution
func BenchmarkRunFindCommand(b *testing.B) {
	// Mock the instance.FindInstances function
	originalFindInstances := instance.FindInstances
	defer func() {
		instance.FindInstances = originalFindInstances
	}()

	instance.FindInstances = func(appCtx *app.ApplicationContext, pattern string, useJSON bool, showDetails bool) error {
		return nil
	}

	appCtx := createMockAppContext()
	cmd := NewFindCmd(appCtx)
	instaceFlags.AllInfoFlag.Add(cmd)
	args := []string{"benchmark-pattern"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunFindCommand(cmd, args, appCtx)
	}
}
