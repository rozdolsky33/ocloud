package database

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewDatabaseCmd(appCtx)

	assert.Equal(t, "database", cmd.Use)
	assert.Contains(t, cmd.Aliases, "db")
	assert.Equal(t, "Explore OCI Database services", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify autonomous subcommand exists
	hasAutonomous := false
	hasHeatWave := false
	for _, sc := range cmd.Commands() {
		if sc.Use == "autonomous" {
			hasAutonomous = true
		}
		if sc.Use == "heatwave" {
			hasHeatWave = true
		}
	}
	assert.True(t, hasAutonomous, "expected autonomous subcommand")
	assert.True(t, hasHeatWave, "expected heatwave subcommand")
}
