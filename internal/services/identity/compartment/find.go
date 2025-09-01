package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/compartment"
)

// FindCompartments searches and displays compartments matching a given name pattern.
func FindCompartments(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	appCtx.Logger.V(logger.Debug).Info("finding compartments", "pattern", namePattern)

	// Create the infrastructure adapter.
	compartmentAdapter := compartment.NewCompartmentAdapter(appCtx.IdentityClient, appCtx.TenancyID)

	// Create the application service, injecting the adapter.
	service := NewService(compartmentAdapter, appCtx.Logger, appCtx.TenancyID)

	ctx := context.Background()
	matchedCompartments, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched compartments: %w", err)
	}

	err = PrintCompartmentsInfo(matchedCompartments, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing matched compartments: %w", err)
	}
	logger.Logger.V(logger.Info).Info("Compartment find operation completed successfully.")
	return nil
}
