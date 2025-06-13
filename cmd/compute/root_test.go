package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestComputeCommand tests the basic structure of the compute command
func TestComputeCommand(t *testing.T) {
	// Create a mock AppContext
	appCtx := &app.AppContext{}

	// Create a new compute command
	cmd := NewComputeCmd(appCtx)

	// Test that the compute command is properly configured
	assert.Equal(t, "compute", cmd.Use)
	assert.Equal(t, "Manage OCI compute resources", cmd.Short)
	assert.Equal(t, "Manage Oracle Cloud Infrastructure compute resources such as instances, images, and more.", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the instance subcommand is added
	instanceCmd := cmd.Commands()
	assert.Equal(t, 1, len(instanceCmd), "compute command should have 1 subcommand")
	assert.Equal(t, "instance", instanceCmd[0].Name(), "compute command should have instance subcommand")
}
