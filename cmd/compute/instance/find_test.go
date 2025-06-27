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

	// Create a new find command
	cmd := NewFindCmd(appCtx)

	// Test that the find command is properly configured
	assert.Equal(t, "find [pattern]", cmd.Use)
	assert.Equal(t, []string{"f"}, cmd.Aliases)
	assert.Equal(t, "Find instances by name pattern", cmd.Short)
	assert.Equal(t, findLong, cmd.Long)
	assert.Equal(t, findExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	// Verify that the command requires exactly one argument
	assert.NotNil(t, cmd.Args)

	// Test that the all flag is added (used for image details)
	imageDetailsFlag := cmd.Flag("all")
	assert.NotNil(t, imageDetailsFlag, "find command should have all flag (used for image details)")
	assert.Equal(t, "all", imageDetailsFlag.Name)
	assert.Equal(t, "A", imageDetailsFlag.Shorthand)

	// Note: The JSON flag is a global flag and is not directly added to the command in the find.go file.
	// It's added at a higher level in the command hierarchy.
}
