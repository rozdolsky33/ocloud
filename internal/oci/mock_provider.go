package oci

import (
	"crypto/rsa"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
)

// MockConfigurationProvider is a mock implementation of common.ConfigurationProvider
// for testing purposes. It returns predefined values for all methods.
type MockConfigurationProvider struct {
	tenancyID      string
	userID         string
	keyFingerprint string
	region         string
	passphrase     string
}

// NewMockConfigurationProvider creates a new MockConfigurationProvider with default test values
func NewMockConfigurationProvider() common.ConfigurationProvider {
	return &MockConfigurationProvider{
		tenancyID:      "ocid1.tenancy.oc1..mock-tenancy-id",
		userID:         "ocid1.user.oc1..mock-user-id",
		keyFingerprint: "mock-key-fingerprint",
		region:         "us-ashburn-1",
		passphrase:     "",
	}
}

// TenancyOCID returns the mock tenancy OCID
func (p *MockConfigurationProvider) TenancyOCID() (string, error) {
	return p.tenancyID, nil
}

// UserOCID returns the mock user OCID
func (p *MockConfigurationProvider) UserOCID() (string, error) {
	return p.userID, nil
}

// KeyFingerprint returns the mock key fingerprint
func (p *MockConfigurationProvider) KeyFingerprint() (string, error) {
	return p.keyFingerprint, nil
}

// Region returns the mock region
func (p *MockConfigurationProvider) Region() (string, error) {
	return p.region, nil
}

// KeyID returns a formatted key ID using the mock values
func (p *MockConfigurationProvider) KeyID() (string, error) {
	tenancy, err := p.TenancyOCID()
	if err != nil {
		return "", err
	}

	user, err := p.UserOCID()
	if err != nil {
		return "", err
	}

	fingerprint, err := p.KeyFingerprint()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", tenancy, user, fingerprint), nil
}

// PrivateRSAKey returns a mock private key
func (p *MockConfigurationProvider) PrivateRSAKey() (key *rsa.PrivateKey, err error) {
	// For testing purposes, we don't need a real private key
	// This is just a placeholder that returns nil
	return nil, nil
}

// Passphrase returns the mock passphrase
func (p *MockConfigurationProvider) Passphrase() (string, error) {
	return p.passphrase, nil
}

// AuthType returns the auth type and configurations
func (p *MockConfigurationProvider) AuthType() (common.AuthConfig, error) {
	return common.AuthConfig{
		AuthType:         common.UserPrincipal,
		IsFromConfigFile: false,
	}, nil
}
