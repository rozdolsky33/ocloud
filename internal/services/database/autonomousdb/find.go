package autonomousdb

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ocidbadapter "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
)

// FindAutonomousDatabases searches for Autonomous Databases matching the provided name pattern in the application context.
// Logs database discovery tasks and can format the result based on the useJSON flag.
func FindAutonomousDatabases(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Autonomous Databases", "pattern", namePattern)

	adapter, err := ocidbadapter.NewAdapter(appCtx.Provider, appCtx.CompartmentID)
	if err != nil {
		return fmt.Errorf("creating database adapter: %w", err)
	}
	service := NewService(adapter, appCtx)

	ctx := context.Background()
	matchedDatabases, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding autonomous databases: %w", err)
	}

	// Convert to domain type for printing
	domainDbs := make([]domain.AutonomousDatabase, 0, len(matchedDatabases))
	for _, db := range matchedDatabases {
		domainDbs = append(domainDbs, domain.AutonomousDatabase(db))
	}

	if err := PrintAutonomousDbInfo(domainDbs, appCtx, nil, useJSON); err != nil {
		return fmt.Errorf("printing autonomous databases: %w", err)
	}

	return nil
}
