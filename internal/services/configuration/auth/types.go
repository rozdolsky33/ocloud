package auth

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/rozdolsky33/ocloud/internal/app"
)

// ConfigProvider is the OCI configuration provider used for authentication.
var ConfigProvider common.ConfigurationProvider

// RegionInfo represents an OCI region with its ID and name.
type RegionInfo struct {
	ID   string
	Name string
}

// AuthenticationResult contains the result of the authentication process.
type AuthenticationResult struct {
	TenancyID   string
	TenancyName string
	Profile     string
	Region      string
}

// Service provides methods for authenticating with OCI.
type Service struct {
	appCtx *app.ApplicationContext
	logger logr.Logger
}
