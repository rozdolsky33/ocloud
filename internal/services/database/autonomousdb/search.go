package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociadb "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
)

// SearchAutonomousDatabases searches for OCI Autonomous Databases matching the given query string in the current context.
// Parameters:
// - appCtx: The application context containing OCI configuration and runtime settings.
// - search: The search query string to perform a fuzzy match against the available databases.
// - useJSON: A flag that determines if the output should be JSON formatted.
// - showAll: A flag indicating whether detailed or summary information should be displayed.
// Returns an error if database search or result processing fails.
func SearchAutonomousDatabases(appCtx *app.ApplicationContext, search string, useJSON bool, showAll bool) error {
	adapter, err := ociadb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating database adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	matchedDatabases, err := service.FuzzySearch(ctx, search)
	if err != nil {
		return fmt.Errorf("finding autonomous databases: %w", err)
	}
	err = PrintAutonomousDbsInfo(matchedDatabases, appCtx, nil, useJSON, showAll)
	if err != nil {
		return fmt.Errorf("printing autonomous databases: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching autonomous databases", "search", search, "matched", len(matchedDatabases))
	return nil
}
