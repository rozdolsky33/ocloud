package heatwave

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
	assert.Equal(t, "Get all HeatWave Databases", cmd.Short)
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

	allFlag := cmd.Flag("all")
	assert.NotNil(t, allFlag, "get command should have all flag")
	assert.Equal(t, "all", allFlag.Name)
	assert.Equal(t, "A", allFlag.Shorthand)
}
