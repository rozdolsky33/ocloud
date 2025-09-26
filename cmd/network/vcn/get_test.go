package vcn

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewGetCmd(appCtx)

	assert.Equal(t, "get", cmd.Use)
	assert.Equal(t, "Get VCNs", cmd.Short)
	assert.Equal(t, getLong, cmd.Long)
	assert.Equal(t, getExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Pagination flags
	limit := cmd.Flag("limit")
	assert.NotNil(t, limit)
	assert.Equal(t, "m", limit.Shorthand)

	page := cmd.Flag("page")
	assert.NotNil(t, page)
	assert.Equal(t, "p", page.Shorthand)

	// Network resource toggles with short aliases
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
