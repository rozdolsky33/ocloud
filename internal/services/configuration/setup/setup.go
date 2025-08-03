package setup

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func SetupTenancyMapping() error {

	s := NewService()
	logger.LogWithLevel(s.logger, 1, "SetupTenancyMapping")
	err := s.ConfigureTenancyFile()
	if err != nil {
		return fmt.Errorf("configuring tenancy mapping file: %w", err)
	}

	return nil
}
