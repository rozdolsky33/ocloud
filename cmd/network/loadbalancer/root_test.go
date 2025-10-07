package loadbalancer

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewLoadBalancerCmd(appCtx)

	assert.Equal(t, "loadbalancer", cmd.Use)
	assert.Contains(t, cmd.Aliases, "lb")
	assert.Contains(t, cmd.Aliases, "lbr")
	assert.Equal(t, "Manage OCI Network Load Balancers", cmd.Short)
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
	assert.True(t, hasSearch, "expected search/find subcommand")
}
