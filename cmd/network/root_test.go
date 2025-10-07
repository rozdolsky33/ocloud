package network

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestNetworkRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewNetworkCmd(appCtx)

	assert.Equal(t, "network", cmd.Use)
	assert.Contains(t, cmd.Aliases, "net")
	assert.Equal(t, "Manage OCI network services", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify subcommands are registered
	hasSubnet := false
	hasVcn := false
	hasLB := false
	for _, sc := range cmd.Commands() {
		switch sc.Use {
		case "subnet":
			hasSubnet = true
		case "vcn":
			hasVcn = true
		case "load-balancer":
			hasLB = true
		}
	}
	assert.True(t, hasSubnet, "expected subnet subcommand")
	assert.True(t, hasVcn, "expected vcn subcommand")
	assert.True(t, hasLB, "expected load-balancer subcommand")
}
