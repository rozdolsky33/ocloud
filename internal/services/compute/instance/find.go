package instance

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindInstances searches for instances in the OCI compartment matching the given name pattern.
// It uses the pre-initialized compute and network clients from the AppContext struct.
// Parameters:
// - appCtx: The application with all clients, logger, and resolved IDs.
// - namePattern: The pattern used to match instance names.
// - showImageDetails: A flag indicating whether to include image details in the output.
// - useJSON: A flag indicating whether to output information in JSON format.
// Returns an error if the operation fails.
func FindInstances(appCtx *app.AppContext, namePattern string, showImageDetails bool, useJSON bool) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(appCtx.Logger, 1, "FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails, "json", useJSON)

	service, err := NewService(appCtx.Provider, appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	matchedInstances, err := service.Find(ctx, namePattern)
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

	// If showImageDetails is true, fetch and display image information
	if showImageDetails {
		// This would be implemented in a future update
		fmt.Println("Image details functionality not yet implemented")
	}

	PrintInstancesTable(matchedInstances, appCtx, nil, useJSON)
	return nil
}
