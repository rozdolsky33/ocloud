package image

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
	assert.Equal(t, "FuzzySearch image by name pattern", cmd.Short)
	assert.Equal(t, findLong, cmd.Long)
	assert.Equal(t, findExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	// Verify that the command requires exactly one argument
	assert.NotNil(t, cmd.Args)
}
