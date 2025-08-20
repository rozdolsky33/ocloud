package subnet

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
	assert.Equal(t, "List all Subnets in the specified tenancy or compartment", cmd.Short)
	assert.Equal(t, listLong, cmd.Long)
	assert.Equal(t, listExamples, cmd.Example)
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

	sortFlag := cmd.Flag("sort")
	assert.NotNil(t, sortFlag, "list command should have sort flag")
	assert.Equal(t, "sort", sortFlag.Name)
	assert.Equal(t, "s", sortFlag.Shorthand)

	// Note: The JSON flag is a global flag and is not directly added to the command in the get.go file.
	// It's added at a higher level in the command hierarchy.
}
