package heatwave

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestSearchCommand tests the basic structure of the search command
func TestSearchCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new search command
	cmd := NewSearchCmd(appCtx)

	// Test that the search command is properly configured
	assert.Equal(t, "search [pattern]", cmd.Use)
	assert.Equal(t, []string{"s"}, cmd.Aliases)
	assert.Equal(t, "Fuzzy Search for HeatWave Databases", cmd.Short)
	assert.Equal(t, searchLong, cmd.Long)
	assert.Equal(t, searchExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	// Verify that the command requires exactly one argument
	assert.NotNil(t, cmd.Args)

	// Test that the --all flags are added
	allFlag := cmd.Flag("all")
	assert.NotNil(t, allFlag, "search command should have all flag")
	assert.Equal(t, "all", allFlag.Name)
	assert.Equal(t, "A", allFlag.Shorthand)

	// Note: The JSON flag is a global flag and is not directly added to the command in the search.go file.
	// It's added at a higher level in the command hierarchy.
}
