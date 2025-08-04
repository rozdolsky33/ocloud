package info

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// ViewConfiguration displays the tenancy mapping information.
// It reads the tenancy-map.yaml file and displays its contents.
// If the realm is not empty, it filters the mappings by the specified realm.
func ViewConfiguration(useJSON bool, realm string) error {
	s := NewService()
	logger.LogWithLevel(s.logger, 1, "ViewConfiguration", "realm", realm)

	result, err := s.LoadTenancyMappings(realm)
	if err != nil {
		return fmt.Errorf("loading tenancy mappings: %w", err)
	}

	err = PrintMappingsFile(result.Mappings, useJSON)
	if err != nil {
		return fmt.Errorf("printing tenancy mappings: %w", err)
	}

	return nil
}
