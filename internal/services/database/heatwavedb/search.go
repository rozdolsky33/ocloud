package heatwavedb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociheatwave "github.com/rozdolsky33/ocloud/internal/oci/database/heatwavedb"
)

// SearchHeatWaveDatabases searches for OCI HeatWave Databases matching the given query string in the current context.
func SearchHeatWaveDatabases(appCtx *app.ApplicationContext, search string, useJSON bool, showAll bool) error {
	adapter, err := ociheatwave.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating HeatWave database adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	matchedDatabases, err := service.FuzzySearch(ctx, search)
	if err != nil {
		return fmt.Errorf("finding HeatWave databases: %w", err)
	}
	err = PrintHeatWaveDbsInfo(matchedDatabases, appCtx, nil, useJSON, showAll)
	if err != nil {
		return fmt.Errorf("printing HeatWave databases: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching HeatWave databases", "search", search, "matched", len(matchedDatabases))
	return nil
}
