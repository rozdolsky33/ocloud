package oci

import (
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// NewComputeClient creates a new OCI compute client using the provided configuration provider.
func NewComputeClient(provider common.ConfigurationProvider) (core.ComputeClient, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating compute client: %w", err)
	}
	return client, nil
}

// NewNetworkClient creates a new OCI virtual network client using the provided configuration provider.
func NewNetworkClient(provider common.ConfigurationProvider) (core.VirtualNetworkClient, error) {
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating virtual network client: %w", err)
	}
	return client, nil
}
