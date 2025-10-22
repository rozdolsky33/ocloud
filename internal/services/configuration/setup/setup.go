package setup

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/logger"
)

// SetupTenancyMapping initializes and configures the tenancy mapping file by using the service's ConfigureTenancyFile method.
// It logs the operation and returns an error if the configuration process fails.
// The context allows for graceful cancellation via Ctrl+C.
func SetupTenancyMapping(ctx context.Context) error {
	s := NewService()
	logger.LogWithLevel(s.logger, logger.Debug, "SetupTenancyMapping")
	err := s.ConfigureTenancyFile(ctx)
	if err != nil {
		if errors.Is(err, ErrCancelled) {
			logger.Logger.V(logger.Info).Info("Tenancy mapping setup cancelled by user.")
			return nil // Don't return an error for user cancellation
		}
		return fmt.Errorf("configuring tenancy mapping file: %w", err)
	}
	logger.Logger.V(logger.Info).Info("Tenancy mapping setup completed successfully.")
	return nil
}
