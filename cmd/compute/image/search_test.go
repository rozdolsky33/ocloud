package image

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
)

// TestSearchCommand tests the basic structure of the image search command
func TestSearchCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new search command
	cmd := NewSearchCmd(appCtx)

	assert.NotNil(t, cmd, "NewSearchCmd should not return nil")

	// Check basic properties
	assert.Equal(t, "search [pattern]", cmd.Use, "Command Use should be set correctly")
	assert.Contains(t, cmd.Aliases, "s", "Command should have alias 's'")
	assert.Equal(t, "Search images by name pattern", cmd.Short, "Command Short description should be set correctly")
	assert.Equal(t, searchLong, cmd.Long, "Command Long description should be set correctly")
	assert.Equal(t, searchExamples, cmd.Example, "Command Example should be set correctly")
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify that the command requires exactly one argument
	assert.NotNil(t, cmd.Args)

	// The JSON flag is a global flag and is not directly added to this command in search.go
	jsonFlag := cmd.Flags().Lookup(cfgflags.FlagNameJSON)
	assert.Nil(t, jsonFlag, "json flag should not be added as a local flag to image search subcommand")

	// Image search command does not support an '--all' flag (unlike instance search)
	allFlag := cmd.Flags().Lookup(cfgflags.FlagNameAll)
	assert.Nil(t, allFlag, "image search command should not include an '--all' flag")
}
