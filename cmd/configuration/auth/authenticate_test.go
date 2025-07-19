package auth

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticateCmd(t *testing.T) {
	// Create a new application context
	appCtx := &app.ApplicationContext{}

	// Create a new authenticate command
	cmd := NewAuthenticateCmd(appCtx)

	// Verify the command properties
	assert.Equal(t, "authenticate", cmd.Use)
	assert.Equal(t, []string{"auth", "a"}, cmd.Aliases)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.NotNil(t, cmd.RunE)

	// Verify the flags
	envOnlyFlag := cmd.Flags().Lookup("env-only")
	assert.NotNil(t, envOnlyFlag)
	assert.Equal(t, "e", envOnlyFlag.Shorthand)
	assert.Equal(t, "false", envOnlyFlag.DefValue)

	filterFlag := cmd.Flags().Lookup("filter")
	assert.NotNil(t, filterFlag)
	assert.Equal(t, "f", filterFlag.Shorthand)
	assert.Equal(t, "", filterFlag.DefValue)
}
