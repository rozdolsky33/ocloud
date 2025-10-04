package oci

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/database"
)

var (
	sharedTransportOnce sync.Once
	sharedTransport     http.RoundTripper
)

// getSharedTransport returns the shared transport.
func getSharedTransport() http.RoundTripper {
	sharedTransportOnce.Do(func() {
		sharedTransport = &http.Transport{
			MaxIdleConns:          256,
			MaxIdleConnsPerHost:   64,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	})
	return sharedTransport
}

// applySharedTransport applies the shared transport to the provided BaseClient.
func applySharedTransport(base *common.BaseClient) {
	if base != nil {
		base.HTTPClient = &http.Client{Transport: getSharedTransport()}
	}
}

// NewIdentityClient creates and returns a new instance of IdentityClient using the provided configuration provider.
func NewIdentityClient(provider common.ConfigurationProvider) (identity.IdentityClient, error) {
	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating identity client: %w", err)
	}
	applySharedTransport(&client.BaseClient)
	return client, nil
}

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

// NewContainerEngineClient creates a new instance of ContainerEngineClient using the provided configuration provider.
func NewContainerEngineClient(provider common.ConfigurationProvider) (containerengine.ContainerEngineClient, error) {
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating container engine client: %w", err)
	}
	return client, nil
}

// NewDatabaseClient creates and returns a new DatabaseClient using the provided configuration.
func NewDatabaseClient(provider common.ConfigurationProvider) (database.DatabaseClient, error) {
	client, err := database.NewDatabaseClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating database client: %w", err)
	}
	return client, nil
}

// NewBastionClient creates and returns a new BastionClient using the specified ConfigurationProvider.
func NewBastionClient(provider common.ConfigurationProvider) (bastion.BastionClient, error) {
	client, err := bastion.NewBastionClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating bastion client: %w", err)
	}
	return client, nil
}

// NewLoadBalancerClient creates and returns a new LoadBalancerClient using the provided configuration provider.
func NewLoadBalancerClient(provider common.ConfigurationProvider) (loadbalancer.LoadBalancerClient, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating load balancer client: %w", err)
	}
	return client, nil
}

// NewCertificatesManagementClient creates and returns a new CertificatesManagementClient.
func NewCertificatesManagementClient(provider common.ConfigurationProvider) (certificatesmanagement.CertificatesManagementClient, error) {
	client, err := certificatesmanagement.NewCertificatesManagementClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating certificates management client: %w", err)
	}
	return client, nil
}

func NewObjectStorageClient(provider common.ConfigurationProvider) (objectstorage.ObjectStorageClient, error) {
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		return client, fmt.Errorf("creating object storage client: %w", err)
	}
	return client, nil
}
