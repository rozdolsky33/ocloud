package objectstorage

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewObjectStorageCmd(appCtx)

	assert.Equal(t, "object-storage", cmd.Use)
	assert.Contains(t, cmd.Aliases, "objectstorage")
	assert.Contains(t, cmd.Aliases, "os")
	assert.Equal(t, "Manage OCI Object Storage: list, get, and search", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify subcommands
	hasGet := false
	hasList := false
	hasSearch := false
	for _, sc := range cmd.Commands() {
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
