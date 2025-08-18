package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindCompartments searches and displays compartments matching a given name pattern with optional JSON formatting.
// It initializes necessary services, performs a fuzzy search, and outputs results using the defined application context.
func FindCompartments(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Compartments", "pattern", namePattern)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compartment service: %w", err)
	}

	ctx := context.Background()
	matchedCompartments, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched compartments: %w", err)
	}

	err = PrintCompartmentsInfo(matchedCompartments, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing matched compartments: %w", err)
	}

	return nil
}
