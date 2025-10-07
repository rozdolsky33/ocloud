package loadbalancer

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewLoadBalancerCmd(appCtx)

	assert.Equal(t, "load-balancer", cmd.Use)
	assert.Contains(t, cmd.Aliases, "loadbalancer")
	assert.Contains(t, cmd.Aliases, "lb")
	assert.Contains(t, cmd.Aliases, "lbr")
	assert.Equal(t, "Manage OCI Network Load Balancers", cmd.Short)
	assert.Equal(t, "Manage Oracle Cloud Infrastructure Network Load Balancers such as LBs, listeners, backend sets, and more", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Sub-commands should be present
	subs := cmd.Commands()
	var hasGet, hasList, hasSearch bool
	for _, sc := range subs {
		switch sc.Use {
		case "get":
			hasGet = true
		case "list":
			hasList = true
		case "search <pattern>":
			hasSearch = true
		}
	}
	assert.True(t, hasGet, "expected get subcommand")
	assert.True(t, hasList, "expected list subcommand")
	assert.True(t, hasSearch, "expected search subcommand")
}
