package instance

import (
	"bytes"
	"strings"
	"testing"

	instaceFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestApplicationContext() *app.ApplicationContext {
	return &app.ApplicationContext{
		CompartmentName: "test-compartment",
	}
}

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	cmd.SetOut(buffer)
	cmd.SetErr(buffer)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buffer.String(), err
}

func TestNewGetCmdMetadata(t *testing.T) {
	cmd := NewGetCmd(newTestApplicationContext())

	assert.Equal(t, "get", cmd.Use, "command use should be set")
	assert.Equal(t, "Paginated Instance Results", cmd.Short, "short description should match")
	assert.Equal(t, getLong, cmd.Long, "long description should match constant")
	assert.Equal(t, getExamples, cmd.Example, "examples should match constant")
	assert.True(t, cmd.SilenceErrors, "command should silence errors")
	assert.True(t, cmd.SilenceUsage, "command should silence usage")
	assert.NotNil(t, cmd.RunE, "RunE should be configured")
}

func TestNewGetCmdAddsExpectedFlags(t *testing.T) {
	cmd := NewGetCmd(newTestApplicationContext())

	limitFlag := cmd.Flags().Lookup(flags.FlagNameLimit)
	require.NotNil(t, limitFlag, "limit flag must be registered")
	assert.Equal(t, flags.FlagShortLimit, limitFlag.Shorthand)
	assert.Equal(t, instaceFlags.FlagDefaultLimit, mustGetInt(t, cmd, flags.FlagNameLimit))

	pageFlag := cmd.Flags().Lookup(flags.FlagNamePage)
	require.NotNil(t, pageFlag, "page flag must be registered")
	assert.Equal(t, flags.FlagShortPage, pageFlag.Shorthand)
	assert.Equal(t, instaceFlags.FlagDefaultPage, mustGetInt(t, cmd, flags.FlagNamePage))

	allFlag := cmd.Flags().Lookup(flags.FlagNameAll)
	require.NotNil(t, allFlag, "all flag must be registered")
	assert.Equal(t, flags.FlagShortAll, allFlag.Shorthand)
	allValue, err := cmd.Flags().GetBool(flags.FlagNameAll)
	require.NoError(t, err)
	assert.False(t, allValue, "all flag default should be false")
}

func TestNewGetCmdHelpOutputMentionsKeyInformation(t *testing.T) {
	cmd := NewGetCmd(newTestApplicationContext())

	output, err := executeCommand(cmd, "--help")
	require.NoError(t, err, "help command should not error")

	expectedSnippets := []string{
		"Paginated Instance Results",
		"Get all instances in the specified compartment",
		"--limit",
		"--page",
		"--all",
		"--json",
	}

	for _, snippet := range expectedSnippets {
		assert.Contains(t, output, snippet, "help output missing snippet %q", snippet)
	}
}

func TestGetCommandConstants(t *testing.T) {
	assert.NotEmpty(t, strings.TrimSpace(getLong), "getLong should not be empty")
	assert.Contains(t, getLong, "Get all instances")
	assert.Contains(t, getLong, "pagination")

	assert.NotEmpty(t, strings.TrimSpace(getExamples), "getExamples should not be empty")
	expectedExamples := []string{
		"ocloud compute instance get",
		"--limit 10 --page 2",
		"--all",
		"--json",
	}
	for _, example := range expectedExamples {
		assert.Contains(t, getExamples, example, "getExamples should include %q", example)
	}
}

func TestFlagDefaultsAndOverrides(t *testing.T) {
	cmd := NewGetCmd(newTestApplicationContext())

	assert.Equal(t, instaceFlags.FlagDefaultLimit, mustGetInt(t, cmd, flags.FlagNameLimit))
	assert.Equal(t, instaceFlags.FlagDefaultPage, mustGetInt(t, cmd, flags.FlagNamePage))

	require.NoError(t, cmd.Flags().Set(flags.FlagNameLimit, "42"))
	require.NoError(t, cmd.Flags().Set(flags.FlagNamePage, "3"))
	require.NoError(t, cmd.Flags().Set(flags.FlagNameAll, "true"))

	assert.Equal(t, 42, mustGetInt(t, cmd, flags.FlagNameLimit))
	assert.Equal(t, 3, mustGetInt(t, cmd, flags.FlagNamePage))

	allValue, err := cmd.Flags().GetBool(flags.FlagNameAll)
	require.NoError(t, err)
	assert.True(t, allValue, "all flag should reflect override to true")
}

func TestRunGetCommandFlagExtractionConsistency(t *testing.T) {
	cmd := &cobra.Command{Use: "get"}
	instaceFlags.LimitFlag.Add(cmd)
	instaceFlags.PageFlag.Add(cmd)
	instaceFlags.AllInfoFlag.Add(cmd)
	cmd.Flags().Bool(flags.FlagNameJSON, false, "output in JSON format")

	require.NoError(t, cmd.Flags().Set(flags.FlagNameLimit, "15"))
	require.NoError(t, cmd.Flags().Set(flags.FlagNamePage, "4"))
	require.NoError(t, cmd.Flags().Set(flags.FlagNameAll, "true"))
	require.NoError(t, cmd.Flags().Set(flags.FlagNameJSON, "true"))

	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, instaceFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, instaceFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	includeAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)

	assert.Equal(t, 15, limit, "limit flag should reflect parsed value")
	assert.Equal(t, 4, page, "page flag should reflect parsed value")
	assert.True(t, useJSON, "json flag should reflect parsed value")
	assert.True(t, includeAll, "all flag should reflect parsed value")
}

func mustGetInt(t *testing.T, cmd *cobra.Command, name string) int {
	t.Helper()
	value, err := cmd.Flags().GetInt(name)
	require.NoError(t, err)
	return value
}
