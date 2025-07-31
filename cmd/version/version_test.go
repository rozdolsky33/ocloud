package version

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rozdolsky33/ocloud/buildinfo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestNewVersionCommand tests the NewVersionCommand function
func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()
	
	// Verify the command properties
	assert.Equal(t, "version", cmd.Use)
	assert.Equal(t, "Print the version information", cmd.Short)
	assert.Equal(t, "Print the version, build time, and git commit hash information", cmd.Long)
	assert.NotNil(t, cmd.RunE, "RunE function should be set")
}

// TestRunCommand tests the runCommand method of VersionInfo
func TestRunCommand(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer
	
	// Create a VersionInfo with the buffer as writer
	vi := &VersionInfo{
		writer: &buf,
		cmd:    &cobra.Command{},
	}
	
	// Run the command
	err := vi.runCommand(vi.cmd, []string{})
	
	// Verify there's no error
	assert.NoError(t, err)
	
	// Verify the output contains version information
	output := buf.String()
	assert.Contains(t, output, "Version:")
	assert.Contains(t, output, "Commit:")
	assert.Contains(t, output, "Built:")
}

// TestPrintVersionInfo tests the PrintVersionInfo function
func TestPrintVersionInfo(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer
	
	// Call the function
	PrintVersionInfo(&buf)
	
	// Verify the output contains version information
	output := buf.String()
	assert.Contains(t, output, "Version:    "+buildinfo.Version)
	assert.Contains(t, output, "Commit:     "+buildinfo.CommitHash)
	assert.Contains(t, output, "Built:      "+buildinfo.BuildTime)
}

// TestPrintVersion tests the PrintVersion function
// This is a simple test that verifies the function doesn't panic
func TestPrintVersion(t *testing.T) {
	// This should not panic
	assert.NotPanics(t, func() {
		PrintVersion()
	}, "PrintVersion should not panic")
}

// TestAddVersionFlag tests the AddVersionFlag function
func TestAddVersionFlag(t *testing.T) {
	// Create a root command
	rootCmd := &cobra.Command{
		Use:   "ocloud",
		Short: "Test command",
	}
	
	// Add the version flag
	AddVersionFlag(rootCmd)
	
	// Verify the flag is added
	flag := rootCmd.PersistentFlags().Lookup("version")
	assert.NotNil(t, flag, "version flag should be added")
	assert.Equal(t, "bool", flag.Value.Type())
	
	// Verify the short flag is added
	shortFlag := rootCmd.PersistentFlags().ShorthandLookup("v")
	assert.NotNil(t, shortFlag, "short version flag should be added")
	assert.Equal(t, flag, shortFlag, "short flag should be the same as the long flag")
	
	// Test the PersistentPreRunE function with version flag
	// Create a buffer to capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	
	// Set the version flag
	rootCmd.PersistentFlags().Set("version", "true")
	
	// Run the PersistentPreRunE function
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	
	// Verify there's no error
	assert.NoError(t, err)
	
	// We can't easily verify the output since PrintVersionInfo writes to os.Stdout
	// In a real test environment, we might use a custom writer that can be captured
}

// TestVersionInfoPrintVersionInfo tests the printVersionInfo method of VersionInfo
func TestVersionInfoPrintVersionInfo(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer
	
	// Create a VersionInfo with the buffer as writer
	vi := &VersionInfo{
		writer: &buf,
	}
	
	// Call the method
	err := vi.printVersionInfo()
	
	// Verify there's no error
	assert.NoError(t, err)
	
	// Verify the output contains version information
	output := buf.String()
	assert.True(t, strings.Contains(output, "Version:"), "Output should contain Version")
	assert.True(t, strings.Contains(output, "Commit:"), "Output should contain Commit")
	assert.True(t, strings.Contains(output, "Built:"), "Output should contain Built")
}