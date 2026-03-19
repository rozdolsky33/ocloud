package dynamicgroup

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestDynamicGroupCommand tests the basic structure of the dynamic-group command
func TestDynamicGroupCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}

	cmd := NewDynamicGroupCmd(appCtx)

	assert.Equal(t, "dynamic-group", cmd.Use)
	assert.Equal(t, []string{"dynamicgroup", "dg"}, cmd.Aliases)
	assert.Equal(t, "Explore OCI Dynamic Groups", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	subCmds := cmd.Commands()
	assert.Equal(t, 3, len(subCmds), "dynamic-group command should have 3 subcommands")

	assert.NotNil(t, findSubCommand(subCmds, "list"))
	assert.NotNil(t, findSubCommand(subCmds, "get"))
	assert.NotNil(t, findSubCommand(subCmds, "search"))
}

func findSubCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
