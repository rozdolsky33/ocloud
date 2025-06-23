package compartment

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func FindCompartments(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Compartments", "pattern", namePattern)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compartment service: %w", err)
	}
	ctx := context.Background()
	compartments, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding compartments: %w", err)
	}

	err = PrintCompartmentsInfo(compartments, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing compartments: %w", err)
	}

	return nil
}
