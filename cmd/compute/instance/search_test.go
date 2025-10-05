package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestFindCommand tests the basic structure of the find command
func TestFindCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new search command
	cmd := NewSearchCmd(appCtx)

	assert.NotNil(t, cmd, "NewSearchCmd should not return nil")

	// Check basic properties
	assert.Equal(t, "search [pattern]", cmd.Use, "Command Use should be set correctly")
	assert.Contains(t, cmd.Aliases, "s", "Command should have alias 's'")
	assert.Equal(t, "Search instances by name pattern", cmd.Short, "Command Short description should be set correctly")
	assert.Equal(t, searchLong, cmd.Long, "Command Long description should be set correctly")
	assert.Equal(t, searchExamples, cmd.Example, "Command Example should be set correctly")
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	// Verify that the command requires exactly one argument
	assert.NotNil(t, cmd.Args)

	// Test that the all flags are added (used for image details)
	imageDetailsFlag := cmd.Flag("all")
	assert.NotNil(t, imageDetailsFlag, "find command should have all flag (used for image details)")
	assert.Equal(t, "all", imageDetailsFlag.Name)
	assert.Equal(t, "A", imageDetailsFlag.Shorthand)

	// Note: The JSON flag is a global flag and is not directly added to the command in the search.go file.
	// It's added at a higher level in the command hierarchy.
}
