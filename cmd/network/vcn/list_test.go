package vcn

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewListCmd(appCtx)

	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "Lists VCNs in a compartment", cmd.Short)
	assert.Equal(t, listLong, cmd.Long)
	assert.Equal(t, listExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	gateway := cmd.Flag("gateway")
	assert.NotNil(t, gateway)
	assert.Equal(t, "G", gateway.Shorthand)

	subnet := cmd.Flag("subnet")
	assert.NotNil(t, subnet)
	assert.Equal(t, "S", subnet.Shorthand)

	nsg := cmd.Flag("nsg")
	assert.NotNil(t, nsg)
	assert.Equal(t, "N", nsg.Shorthand)

	route := cmd.Flag("route-table")
	assert.NotNil(t, route)
	assert.Equal(t, "R", route.Shorthand)

	secList := cmd.Flag("security-list")
	assert.NotNil(t, secList)
	assert.Equal(t, "L", secList.Shorthand)
}
