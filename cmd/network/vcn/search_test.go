package vcn

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestSearchCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewSearchCmd(appCtx)

	assert.Equal(t, "search <pattern>", cmd.Use)
	assert.Equal(t, "Fuzzy search for VCNs", cmd.Short)
	assert.Equal(t, searchLong, cmd.Long)
	assert.Equal(t, searchExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Aliases
	assert.Contains(t, cmd.Aliases, "s")

	// Require exactly one arg
	assert.NotNil(t, cmd.Args)

	// Flags added with short aliases
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
