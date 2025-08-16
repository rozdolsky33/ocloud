package info

import (
	"fmt"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// NewService initializes a new Service instance with the provided application context.
func NewService() *Service {
	appCtx := &app.ApplicationContext{
		Logger: logger.Logger,
	}
	service := &Service{
		logger: appCtx.Logger,
	}
	return service
}

// LoadTenancyMappings loads the tenancy mappings from the file and filters them by realm if specified.
// It returns a TenancyMappingResult containing the filtered mappings and an error if encountered.
func (s *Service) LoadTenancyMappings(realm string) (*TenancyMappingResult, error) {
	// Load the tenancy mapping from the file
	logger.LogWithLevel(s.logger, 3, "Loading tenancy mappings", "realm", realm)
	tenancies, err := appConfig.LoadTenancyMap()
	if err != nil {
		return nil, fmt.Errorf("loading tenancy map: %w", err)
	}

	// Filter by realm if specified
	var filteredMappings []appConfig.MappingsFile
	for _, tenancy := range tenancies {
		if realm != "" && !strings.EqualFold(tenancy.Realm, realm) {
			continue
		}

		filteredMappings = append(filteredMappings, tenancy)
	}

	logger.LogWithLevel(s.logger, 3, "Loaded tenancy mappings", "count", len(filteredMappings))
	return &TenancyMappingResult{
		Mappings: filteredMappings,
	}, nil
}
