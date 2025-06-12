package oci

import (
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"
)

// TestNewComputeClient tests the NewComputeClient function
func TestNewComputeClient(t *testing.T) {
	// Since we can't easily mock the OCI SDK's core.NewComputeClientWithConfigurationProvider function,
	// we'll just test that our function doesn't panic and returns a non-nil client
	// when given a valid configuration provider.

	// Create a real configuration provider (this won't make actual API calls)
	provider := common.DefaultConfigProvider()

	// Call the function
	client, err := NewComputeClient(provider)

	// Check that the function returns a non-nil client and no error
	// Note: This doesn't actually test the function's behavior, just that it doesn't crash
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

// TestNewNetworkClient tests the NewNetworkClient function
func TestNewNetworkClient(t *testing.T) {
	// Since we can't easily mock the OCI SDK's core.NewVirtualNetworkClientWithConfigurationProvider function,
	// we'll just test that our function doesn't panic and returns a non-nil client
	// when given a valid configuration provider.

	// Create a real configuration provider (this won't make actual API calls)
	provider := common.DefaultConfigProvider()

	// Call the function
	client, err := NewNetworkClient(provider)

	// Check that the function returns a non-nil client and no error
	// Note: This doesn't actually test the function's behavior, just that it doesn't crash
	assert.NotNil(t, client)
	assert.NoError(t, err)
}
