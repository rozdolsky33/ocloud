package loadbalancer

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewGetCmd(appCtx)

	assert.Equal(t, "get", cmd.Use)
	assert.Equal(t, " Get Load Balancer Paginated Results", cmd.Short)
	assert.Equal(t, getLong, cmd.Long)
	assert.Equal(t, getExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Flags added
	limit := cmd.Flag("limit")
	assert.NotNil(t, limit)
	assert.Equal(t, "m", limit.Shorthand)

	page := cmd.Flag("page")
	assert.NotNil(t, page)
	assert.Equal(t, "p", page.Shorthand)

	all := cmd.Flag("all")
	assert.NotNil(t, all)
	assert.Equal(t, "A", all.Shorthand)
}
