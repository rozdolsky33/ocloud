package instance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config"
)

// TestNewFindCmd tests the newFindCmd function
func TestNewFindCmd(t *testing.T) {
	// Create a mock AppContext
	mockCtx := &app.AppContext{
		Ctx: context.Background(),
	}

	// Call the function
	cmd := newFindCmd(mockCtx)

	// Test that the command is properly configured
	assert.Equal(t, "find", cmd.Name())
	assert.Equal(t, "Find instances by name pattern", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.NotNil(t, cmd.RunE)
	// Can't directly compare functions, so check that Args is set and behaves like ExactArgs(1)
	assert.NotNil(t, cmd.Args)
	// Test with 0 args (should fail)
	assert.Error(t, cmd.Args(cmd, []string{}))
	// Test with 1 arg (should pass)
	assert.NoError(t, cmd.Args(cmd, []string{"test-pattern"}))
	// Test with 2 args (should fail)
	assert.Error(t, cmd.Args(cmd, []string{"test-pattern", "extra-arg"}))

	// Test that the flags are added
	imageDetailsFlag := cmd.Flags().Lookup(config.FlagNameImageDetails)
	assert.NotNil(t, imageDetailsFlag, "image-details flag should be added")
	assert.Equal(t, config.FlagShortImageDetails, imageDetailsFlag.Shorthand)
	assert.Equal(t, config.FlagDescImageDetails, imageDetailsFlag.Usage)
}
