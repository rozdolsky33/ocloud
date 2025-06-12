package instance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestNewListCmd tests the newListCmd function
func TestNewListCmd(t *testing.T) {
	// Create a mock AppContext
	mockCtx := &app.AppContext{
		Ctx: context.Background(),
	}

	// Call the function
	cmd := newListCmd(mockCtx)

	// Test that the command is properly configured
	assert.Equal(t, "list", cmd.Name())
	assert.Equal(t, "List all instances in the compartment", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.NotNil(t, cmd.RunE)
}
