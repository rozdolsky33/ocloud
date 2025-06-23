package compartment

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func ListCompartments(appCtx *app.ApplicationContext, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing Compartments")

	service, err := NewService(appCtx)

	if err != nil {
		return fmt.Errorf("creating compartment service: %w", err)
	}
	ctx := context.Background()
	compartments, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("listing compartments: %w", err)
	}

	err = PrintCompartmentsInfo(compartments, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing compartments: %w", err)
	}
	return nil
}
