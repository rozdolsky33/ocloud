package loadbalancer

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestFindCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewFindCmd(appCtx)

	assert.Equal(t, "find <pattern>", cmd.Use)
	assert.Contains(t, cmd.Aliases, "f")
	assert.Equal(t, "Finds Load Balancer with existing attribute", cmd.Short)
	assert.Equal(t, findLong, cmd.Long)
	assert.Equal(t, findExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Require exactly one arg
	assert.NotNil(t, cmd.Args)

	// Flags added
	all := cmd.Flag("all")
	assert.NotNil(t, all)
	assert.Equal(t, "A", all.Shorthand)
}
