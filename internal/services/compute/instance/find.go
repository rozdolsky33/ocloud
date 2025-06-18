package instance

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindInstances searches for instances in the OCI compartment matching the given name pattern.
// It uses the pre-initialized compute and network clients from the ApplicationContext struct.
// Parameters:
// - appCtx: The application with all clients, logger, and resolved IDs.
// - namePattern: The pattern used to match instance names.
// - showImageDetails: A flag indicating whether to include images details in the output.
// - useJSON: A flag indicating whether to output information in JSON format.
// Returns an error if the operation fails.
func FindInstances(appCtx *app.ApplicationContext, namePattern string, showImageDetails bool, useJSON bool) error {
	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(appCtx.Logger, 1, "FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails, "json", useJSON)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	matchedInstances, err := service.Find(ctx, namePattern, showImageDetails)
	if err != nil {
		return fmt.Errorf("finding instances: %w", err)
	}

	// Display matched instances
	if len(matchedInstances) == 0 {
		if useJSON {
			// Return an empty JSON array if no instances found
			fmt.Println(`{"instances": [], "pagination": null}`)
		} else {
			fmt.Printf("No instances found matching pattern: %s\n", namePattern)
		}
		return nil
	}

	// Pass the showImageDetails flag to PrintInstancesInfo
	err = PrintInstancesInfo(matchedInstances, appCtx, nil, useJSON, showImageDetails)
	if err != nil {
		return fmt.Errorf("printing instances table: %w", err)
	}
	return nil
}
