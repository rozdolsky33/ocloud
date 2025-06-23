package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewComputeClient tests the NewComputeClient function
func TestNewComputeClient(t *testing.T) {
	// Use our mock configuration provider instead of the real one
	provider := NewMockConfigurationProvider()

	// Call the function
	client, err := NewComputeClient(provider)

	// Check that the function returns a non-nil client and no error
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

// TestNewNetworkClient tests the NewNetworkClient function
func TestNewNetworkClient(t *testing.T) {
	// Use our mock configuration provider instead of the real one
	provider := NewMockConfigurationProvider()

	// Call the function
	client, err := NewNetworkClient(provider)

	// Check that the function returns a non-nil client and no error
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewContainerEngineClient(t *testing.T) {
	// Use our mock configuration provider instead of the real one
	provider := NewMockConfigurationProvider()

	// Call the function
	client, err := NewContainerEngineClient(provider)

	assert.NotNil(t, client)
	assert.NoError(t, err)
}
