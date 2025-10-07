package loadbalancer

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestSearchCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewSearchCmd(appCtx)

	assert.Equal(t, "search <pattern>", cmd.Use)
	assert.Contains(t, cmd.Aliases, "s")
	assert.Equal(t, "Fuzzy search for Load Balancers", cmd.Short)
	assert.Equal(t, searchLong, cmd.Long)
	assert.Equal(t, searchExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Require exactly one arg
	assert.NotNil(t, cmd.Args)

	// Flags added
	all := cmd.Flag("all")
	assert.NotNil(t, all)
	assert.Equal(t, "A", all.Shorthand)
}
