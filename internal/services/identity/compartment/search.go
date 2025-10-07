package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/compartment"
)

// SearchCompartments searches and displays compartments matching a given name pattern.
func SearchCompartments(appCtx *app.ApplicationContext, namePattern string, useJSON bool, ocid string) error {
	ctx := context.Background()
	compartmentAdapter := compartment.NewCompartmentAdapter(appCtx.IdentityClient, ocid)

	// Create the application service, injecting the adapter.
	service := NewService(compartmentAdapter, appCtx.Logger, ocid)

	matchedCompartments, err := service.FuzzySearch(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched compartments: %w", err)
	}
	err = PrintCompartmentsInfo(matchedCompartments, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing matched compartments: %w", err)
	}

	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching compartments", "search", namePattern, "matched", len(matchedCompartments))
	return nil
}
