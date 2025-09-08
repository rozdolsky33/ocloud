package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociadb "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
)

// FindAutonomousDatabases searches for Autonomous Databases matching the provided name pattern in the application context.
// Logs database discovery tasks and can format the result based on the useJSON flag.
func FindAutonomousDatabases(appCtx *app.ApplicationContext, namePattern string, useJSON bool, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Finding Autonomous Databases", "pattern", namePattern)

	adapter, err := ociadb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating database adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	matchedDatabases, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding autonomous databases: %w", err)
	}

	return PrintAutonomousDbsInfo(matchedDatabases, appCtx, nil, useJSON, showAll)
}
