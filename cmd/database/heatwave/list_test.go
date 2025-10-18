package heatwave

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestListCommand tests the basic structure of the list command
func TestListCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new list command
	cmd := NewListCmd(appCtx)

	// Test that the list command is properly configured
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, []string{"l"}, cmd.Aliases)
	assert.Equal(t, "List all HeatWave Databases", cmd.Short)
	assert.Equal(t, listLong, cmd.Long)
	assert.Equal(t, listExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Note: The JSON flag is a global flag and is not directly added to the command in the list.go file.
	// It's added at a higher level in the command hierarchy.
}
