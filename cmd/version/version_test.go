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
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a VersionInfo with the buffer as a writer
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
	// Create a buffer to capture the output
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

// TestVersionInfoPrintVersionInfo tests the printVersionInfo method of VersionInfo
func TestVersionInfoPrintVersionInfo(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a VersionInfo with the buffer as a writer
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

func TestAddVersionFlag(t *testing.T) {
	// Create a root command
	rootCmd := &cobra.Command{
		Use: "ocloud",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	// Create a buffer to capture the output
	buf := new(bytes.Buffer)

	// Add the version flag
	AddVersionFlag(rootCmd, buf)

	// Verify the flag is added
	flag := rootCmd.PersistentFlags().Lookup("version")
	assert.NotNil(t, flag, "version flag should be added")

	// --- Test case 1: --version flag is set ---
	rootCmd.SetArgs([]string{"--version"})
	err := rootCmd.Execute()
	assert.NoError(t, err, "Execute should not return an error when --version is set")

	// Verify the output contains version information
	output := buf.String()
	assert.Contains(t, output, "Version:")
	assert.Contains(t, output, "Commit:")
	assert.Contains(t, output, "Built:")

	// --- Test case 2: --version flag is NOT set ---
	var originalPreRunCalled bool
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		originalPreRunCalled = true
		return nil
	}

	rootCmd.SetArgs([]string{})
	err = rootCmd.Execute()
	assert.NoError(t, err, "Execute should not return an error when --version is not set")

	// Verify that the original PersistentPreRunE was called
	assert.True(t, originalPreRunCalled, "Original PersistentPreRunE should be called when --version is not set")
}
