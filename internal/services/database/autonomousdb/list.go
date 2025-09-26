package autonomousdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	ociadb "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListAutonomousDatabases lists all Autonomous Databases and presents an interactive TUI for the user to select one.
// It retrieves the selected database and prints its details; if the user cancels the selection the function returns nil.
// It returns an error if listing databases, selecting (except cancellation), retrieving, or printing the database fails.
func ListAutonomousDatabases(appCtx *app.ApplicationContext, useJSON bool) error {
	ctx := context.Background()
	autonomousDatabaseAdapter, err := ociadb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating autonomous database adapter: %w", err)
	}
	service := NewService(autonomousDatabaseAdapter, appCtx)
	allDatabases, err := service.ListAutonomousDb(ctx)

	if err != nil {
		return fmt.Errorf("listing autonomous databases: %w", err)
	}

	//TUI
	model := ociadb.NewDatabaseListModel(allDatabases)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting database: %w", err)
	}

	database, err := service.repo.GetAutonomousDatabase(ctx, id)
	if err != nil {
		return fmt.Errorf("getting database: %w", err)
	}

	return PrintAutonomousDbInfo(database, appCtx, useJSON, true)
}
