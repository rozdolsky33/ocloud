package bastion

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	ocidbadapter "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
	autonomousdbsvc "github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

// connectDatabase runs the DB target flow. We can’t always auto-verify reachability,
// so we surface that limitation to the user.
func connectDatabase(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

	adapter, err := ocidbadapter.NewAdapter(appCtx.Provider, appCtx.CompartmentID)
	if err != nil {
		return fmt.Errorf("error creating database adapter: %w", err)
	}
	dbService := autonomousdbsvc.NewService(adapter, appCtx)

	dbs, _, _, err := dbService.List(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("list databases: %w", err)
	}
	if len(dbs) == 0 {
		fmt.Println("No Autonomous Databases found.")
		return nil
	}

	// Port 1521 or 1522 is the default ports for Oracle Database

	dm := NewDBListModelFancy(dbs)
	dp := tea.NewProgram(dm, tea.WithContext(ctx))
	dres, err := dp.Run()
	if err != nil {
		return fmt.Errorf("DB selection TUI: %w", err)
	}
	chosen, ok := dres.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return ErrAborted
	}

	var db autonomousdbsvc.AutonomousDatabase
	for _, d := range dbs {
		if d.ID == chosen.Choice() {
			db = d
			break
		}
	}

	_, reason := svc.CanReach(ctx, b, "", "")
	fmt.Println("Reachability to DB cannot be automatically verified:", reason)
	fmt.Printf("Selected database: %s (ID: %s)\n", db.Name, db.ID)
	fmt.Printf("\n---\nPrepared %s session on Bastion %s (ID: %s) to database %s.\n",
		sType, b.Name, b.ID, db.Name)
	return nil
}
