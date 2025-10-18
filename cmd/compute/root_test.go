package compute

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestComputeCommand tests the basic structure of the compute command
func TestComputeCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new compute command
	cmd := NewComputeCmd(appCtx)

	// Test that the compute command is properly configured
	assert.Equal(t, "compute", cmd.Use)
	assert.Equal(t, "Explore OCI compute services", cmd.Short)
	assert.Equal(t, "Explore Oracle Cloud Infrastructure Compute services such as instances, images, and oke.", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 3, len(subCmds), "compute command should have 3 subcommands")

	// Check that the instance subcommand is present
	instanceCmd := computeSubCommand(subCmds, "instance")
	assert.NotNil(t, instanceCmd, "compute command should have instance subcommand")

	// Check that the image subcommand is present
	imageCmd := computeSubCommand(subCmds, "image")
	assert.NotNil(t, imageCmd, "compute command should have image subcommand")

	// Check that the oke subcommand is present
	okeCmd := computeSubCommand(subCmds, "oke")
	assert.NotNil(t, okeCmd, "compute command should have oke subcommand")
}

// computeSubCommand is a helper function to find a subcommand by name
func computeSubCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
