package objectstorage

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestSearchCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewSearchCmd(appCtx)

	assert.Equal(t, "search <pattern>", cmd.Use)
	assert.Equal(t, "Fuzzy search for Buckets", cmd.Short)
	assert.Equal(t, searchLong, cmd.Long)
	assert.Equal(t, searchExamples, cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Aliases
	assert.Contains(t, cmd.Aliases, "s")

	// Args required
	assert.NotNil(t, cmd.Args)
}
