package identity

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestIdentityRootCommand(t *testing.T) {
	appCtx := &app.ApplicationContext{}
	cmd := NewIdentityCmd(appCtx)

	assert.Equal(t, "identity", cmd.Use)
	assert.Contains(t, cmd.Aliases, "ident")
	assert.Contains(t, cmd.Aliases, "idt")
	assert.Equal(t, "Explore OCI identity services and manage bastion sessions", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify subcommands are registered
	hasBastion := false
	hasCompartment := false
	hasPolicy := false
	for _, sc := range cmd.Commands() {
		switch sc.Use {
		case "bastion":
			hasBastion = true
		case "compartment":
			hasCompartment = true
		case "policy":
			hasPolicy = true
		}
	}
	assert.True(t, hasBastion, "expected bastion subcommand")
	assert.True(t, hasCompartment, "expected compartment subcommand")
	assert.True(t, hasPolicy, "expected policy subcommand")
}
