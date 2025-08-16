package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindAutonomousDatabases searches for Autonomous Databases matching the provided name pattern in the application context.
// Logs database discovery tasks and can format the result based on the useJSON flag.
// Returns an error if the discovery or result formatting fails.
func FindAutonomousDatabases(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Autonomous Databases", "pattern", namePattern)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating autonomous database service: %w", err)
	}

	ctx := context.Background()
	matchedDatabases, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding autonomous databases: %w", err)
	}

	err = PrintAutonomousDbInfo(matchedDatabases, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing autonomous databases: %w", err)
	}

	return nil
}
