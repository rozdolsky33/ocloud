package bastion

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func ListBastions(ctx context.Context, appCtx *app.ApplicationContext, useJSON bool) error {

	logger.LogWithLevel(appCtx.Logger, 1, "Listing bastions")

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating bastion service: %w", err)
	}

	bastions, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("listing bastions: %w", err)
	}

	err = PrintBastionInfo(bastions, appCtx, useJSON)
	if err != nil {
		return fmt.Errorf("printing bastions: %w", err)
	}

	return nil
}
