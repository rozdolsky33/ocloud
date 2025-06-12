package instance

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestNewListCmd tests the newListCmd function
func TestNewListCmd(t *testing.T) {
	// Create a mock AppContext
	mockApp := &app.AppContext{
		// Initialize with minimal required fields for the test
		Logger: logr.Discard(),
	}

	// Call the function
	cmd := newListCmd(mockApp)

	// Test that the command is properly configured
	assert.Equal(t, "list", cmd.Name())
	assert.Equal(t, "List all instances in the compartment", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.NotNil(t, cmd.RunE)
}
