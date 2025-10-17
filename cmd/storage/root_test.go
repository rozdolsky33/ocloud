package storage

import (
	"strings"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestStorageRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewStorageCmd(appCtx)

	assert.Equal(t, "storage", cmd.Use)
	assert.Contains(t, cmd.Aliases, "stg")
	assert.Equal(t, "Explore OCI Storage Resources", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify object-storage subcommand exists and has expected aliases and subcommands
	var osCmd *cobra.Command
	for _, sc := range cmd.Commands() {
		if sc.Use == "object-storage" {
			osCmd = sc
			break
		}
	}
	if assert.NotNil(t, osCmd, "expected object-storage subcommand") {
		// Aliases
		assert.Contains(t, osCmd.Aliases, "objectstorage")
		assert.Contains(t, osCmd.Aliases, "os")

		// Subcommands
		hasGet, hasList, hasSearch := false, false, false
		for _, sub := range osCmd.Commands() {
			switch sub.Use {
			case "get":
				hasGet = true
			case "list":
				hasList = true
			default:
				if strings.HasPrefix(sub.Use, "search") {
					hasSearch = true
				}
			}
		}
		assert.True(t, hasGet, "expected get subcommand")
		assert.True(t, hasList, "expected list subcommand")
		assert.True(t, hasSearch, "expected search subcommand")
	}
}
