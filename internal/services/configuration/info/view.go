package info

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// ViewConfiguration displays the tenancy mapping information.
// It reads the tenancy-map.yaml file and displays its contents.
// If the realm is not empty, it filters the mappings by the specified realm.
func ViewConfiguration(appCtx *app.ApplicationContext, useJSON bool, realm string) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Viewing tenancy mapping configuration", "realm", realm)

	// Create a new service
	service := NewService(appCtx)

	// Load tenancy mappings
	result, err := service.LoadTenancyMappings(realm)
	if err != nil {
		return fmt.Errorf("loading tenancy mappings: %w", err)
	}

	// Display tenancy mapping information
	err = PrintMappingsFile(result.Mappings, appCtx, useJSON)
	if err != nil {
		return fmt.Errorf("printing tenancy mappings: %w", err)
	}

	return nil
}
