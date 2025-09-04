package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetupMappingFile tests the SetupMappingFile function
func TestSetupMappingFile(t *testing.T) {
	// Call the function
	cmd := SetupMappingFile()

	// Verify the command properties
	assert.NotNil(t, cmd, "Command should not be nil")
	assert.Equal(t, "setup", cmd.Use, "Command Use should be 'setup'")
	assert.Equal(t, "Create tenancy mapping file or add a record", cmd.Short, "Command Short description should match")
	assert.NotNil(t, cmd.RunE, "Command RunE function should not be nil")
}
