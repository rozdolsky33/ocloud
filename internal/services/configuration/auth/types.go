package auth

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
)

// result holds the result of the authentication process, including tenancy, compartment, profile, and region details.
var result *AuthenticationResult
var err error

// RegionInfo represents an OCI region with its ID and name.
type RegionInfo struct {
	ID   string
	Name string
}

// AuthenticationResult contains the result of the authentication process.
type AuthenticationResult struct {
	TenancyID       string
	TenancyName     string
	CompartmentName string
	Profile         string
	Region          string
}

// Service represents an authentication service using OCI configuration and logging utilities.
type Service struct {
	Provider common.ConfigurationProvider
	logger   logr.Logger
}
