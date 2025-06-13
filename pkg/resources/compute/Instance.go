package compute

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// ListInstances lists all instances in the configured compartment using the provided application.
// It uses the pre-initialized compute client from the AppContext struct.
func ListInstances(application *app.AppContext) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(application.Logger, 1, "ListInstances()")

	// Use the pre-initialized compute client from the AppContext struct
	// No need to create a new client

	fmt.Println("Inside Instance resources running List Instances")
	// Use application.ComputeClient to list instances
	// ...

	return nil
}

// FindInstances searches for instances in the OCI compartment matching the given name pattern.
// It uses the pre-initialized compute and network clients from the AppContext struct.
// Parameters:
// - application: The application with all clients, logger, and resolved IDs.
// - namePattern: The pattern used to match instance names.
// - showImageDetails: A flag indicating whether to include image details in the output.
// Returns an error if the operation fails.
func FindInstances(application *app.AppContext, namePattern string, showImageDetails bool) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(application.Logger, 1, "FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails)

	// Use the pre-initialized compute and network clients from the AppContext struct
	// No need to create new clients

	// Use application.ComputeClient and application.NetworkClient to find instances
	// ...

	return nil
}
