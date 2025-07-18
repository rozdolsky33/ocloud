package info

import (
	"github.com/go-logr/logr"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
)

// Service provides operations and functionalities related to tenancy mapping information.
type Service struct {
	logger logr.Logger
}

// TenancyMappingResult represents the result of loading and filtering tenancy mappings.
type TenancyMappingResult struct {
	Mappings []appConfig.MappingsFile
}
