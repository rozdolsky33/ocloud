package vcn

import (
	"bytes"
	"github.com/go-logr/logr/testr"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewVcnCmd(t *testing.T) {
	appCtx := &app.ApplicationContext{
		Logger: testr.New(t),
	}

	cmd := NewVcnCmd(appCtx)

	assert.NotNil(t, cmd)
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	assert.Equal(t, "vcn", cmd.Use)
	assert.Equal(t, "Manage OCI Virtual Cloud Networks (VCNs)", cmd.Short)

	// Check if subcommands are added
	expectedSubcommands := []string{"get", "list"}
	for _, sub := range expectedSubcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == sub {
				found = true
				break
			}
		}
		assert.True(t, found, "subcommand %s not found", sub)
	}
}
