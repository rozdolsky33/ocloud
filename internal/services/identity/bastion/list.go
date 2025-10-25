package bastion

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// GetBastions retrieves a list of bastion hosts and displays their information, optionally in JSON format.
func GetBastions(ctx context.Context, appCtx *app.ApplicationContext, useJSON bool) error {

	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing bastions")

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating bastion service: %w", err)
	}

	bastions, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("listing bastions: %w", err)
	}

	return PrintBastionInfo(bastions, appCtx, useJSON)
}
