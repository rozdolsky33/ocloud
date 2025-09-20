package compartment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestGetCommand tests the basic structure of the get command
func TestGetCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new get command
	cmd := NewGetCmd(appCtx)

	// Test that the get command is properly configured
	assert.Equal(t, "get", cmd.Use)
	assert.Equal(t, "Get all Compartments in the specified tenancy or compartment", cmd.Short)
	assert.Equal(t, getLong, cmd.Long)
	assert.Equal(t, getExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the flags are added
	limitFlag := cmd.Flag("limit")
	assert.NotNil(t, limitFlag, "get command should have limit flag")
	assert.Equal(t, "limit", limitFlag.Name)
	assert.Equal(t, "m", limitFlag.Shorthand)

	pageFlag := cmd.Flag("page")
	assert.NotNil(t, pageFlag, "get command should have page flag")
	assert.Equal(t, "page", pageFlag.Name)
	assert.Equal(t, "p", pageFlag.Shorthand)

	// Scope-related flags
	scopeFlag := cmd.Flag("scope")
	assert.NotNil(t, scopeFlag, "get command should have scope flag")

	tenancyScopeFlag := cmd.Flag("tenancy-scope")
	assert.NotNil(t, tenancyScopeFlag, "get command should have tenancy-scope flag")
	assert.Equal(t, "T", tenancyScopeFlag.Shorthand)

	// Note: The JSON flag is a global flag and is not directly added to the command in the gateway.go file.
	// It's added at a higher level in the command hierarchy.
}
