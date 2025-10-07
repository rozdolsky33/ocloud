package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestListCommand tests the basic structure of the list command
func TestGetCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new list command
	cmd := NewGetCmd(appCtx)

	// Test that the list command is properly configured
	assert.Equal(t, "get", cmd.Use)
	assert.Equal(t, "Paginated Instance Results", cmd.Short)
	assert.Equal(t, getLong, cmd.Long)
	assert.Equal(t, getExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the flags are added
	limitFlag := cmd.Flag("limit")
	assert.NotNil(t, limitFlag, "list command should have limit flag")
	assert.Equal(t, "limit", limitFlag.Name)
	assert.Equal(t, "m", limitFlag.Shorthand)

	pageFlag := cmd.Flag("page")
	assert.NotNil(t, pageFlag, "list command should have page flag")
	assert.Equal(t, "page", pageFlag.Name)
	assert.Equal(t, "p", pageFlag.Shorthand)

	// Test that the all flags are added (used for image details)
	imageDetailsFlag := cmd.Flag("all")
	assert.NotNil(t, imageDetailsFlag, "list command should have all flag (used for image details)")
	assert.Equal(t, "all", imageDetailsFlag.Name)
	assert.Equal(t, "A", imageDetailsFlag.Shorthand)
}
