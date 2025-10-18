package heatwavedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	ociheatwave "github.com/rozdolsky33/ocloud/internal/oci/database/heatwavedb"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListHeatWaveDatabases lists all HeatWave Databases in the application context with TUI.
func ListHeatWaveDatabases(appCtx *app.ApplicationContext, useJSON bool) error {
	ctx := context.Background()
	heatwaveDatabaseAdapter, err := ociheatwave.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating HeatWave database adapter: %w", err)
	}
	service := NewService(heatwaveDatabaseAdapter, appCtx)
	allDatabases, err := service.ListHeatWaveDb(ctx)

	if err != nil {
		return fmt.Errorf("listing HeatWave databases: %w", err)
	}

	// TUI
	model := ociheatwave.NewDatabaseListModel(allDatabases)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting database: %w", err)
	}

	database, err := service.repo.GetHeatWaveDatabase(ctx, id)
	if err != nil {
		return fmt.Errorf("getting database: %w", err)
	}

	return PrintHeatWaveDbInfo(database, appCtx, useJSON, true)
}
