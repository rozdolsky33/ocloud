package setup

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/logger"
)

// SetupTenancyMapping initializes and configures the tenancy mapping file by using the service's ConfigureTenancyFile method.
// It logs the operation and returns an error if the configuration process fails.
func SetupTenancyMapping() error {
	s := NewService()
	logger.LogWithLevel(s.logger, logger.Debug, "SetupTenancyMapping")
	err := s.ConfigureTenancyFile()
	if err != nil {
		return fmt.Errorf("configuring tenancy mapping file: %w", err)
	}
	logger.Logger.V(logger.Info).Info("Tenancy mapping setup completed successfully.")
	return nil
}
