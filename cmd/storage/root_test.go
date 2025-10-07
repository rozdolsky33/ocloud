package storage

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestStorageRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewStorageCmd(appCtx)

	assert.Equal(t, "storage", cmd.Use)
	assert.Contains(t, cmd.Aliases, "stg")
	assert.Equal(t, "Manage OCI Storage Resources", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify object-storage subcommand exists
	hasOS := false
	for _, sc := range cmd.Commands() {
		if sc.Use == "object-storage" {
			hasOS = true
			break
		}
	}
	assert.True(t, hasOS, "expected object-storage subcommand")
}
