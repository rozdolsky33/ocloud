package auth

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
)

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

// Service represents a service for handling OCI configuration and authentication processes.
type Service struct {
	Provider common.ConfigurationProvider
	logger   logr.Logger
}
