package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociadb "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
)

// SearchAutonomousDatabases searches for Autonomous Databases matching the provided name pattern in the application context.
// Logs database discovery tasks and can format the result based on the useJSON flag.
func SearchAutonomousDatabases(appCtx *app.ApplicationContext, namePattern string, useJSON bool, showAll bool) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Finding Autonomous Databases", "pattern", namePattern)

	adapter, err := ociadb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating database adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	matchedDatabases, err := service.FuzzySearch(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding autonomous databases: %w", err)
	}
	err = PrintAutonomousDbsInfo(matchedDatabases, appCtx, nil, useJSON, true)
	if err != nil {
		return fmt.Errorf("printing autonomous databases: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching autonomous databases", "searchPattern", namePattern, "matched", len(matchedDatabases))
	return nil
}
