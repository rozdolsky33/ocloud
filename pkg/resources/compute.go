package resources

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// ListInstances lists all instances in the configured compartment using the provided application context.
// It initializes a compute client based on the OCI configuration and logs the execution process.
func ListInstances(appCtx *app.AppContext) error {
	logger.Logger.V(1).Info("ListInstances()")

	// Create a compute client
	_, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	fmt.Println("Inside Instance resources running List Instances")
	// Use computeClient to list instances
	// ...

	return nil
}

// FindInstances searches for instances in the OCI compartment matching the given name pattern.
// It initializes OCI compute and network clients and optionally displays image details.
// Parameters:
// - appCtx: The application context containing OCI configuration and resolved IDs.
// - namePattern: The pattern used to match instance names.
// - showImageDetails: A flag indicating whether to include image details in the output.
// Returns an error if the clients cannot be created or the operation fails.
func FindInstances(appCtx *app.AppContext, namePattern string, showImageDetails bool) error {
	logger.Logger.V(1).Info("FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails)

	// Create a compute client
	_, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	// Create a network client
	_, err = oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	// Use computeClient and networkClient to find instances
	// ...

	return nil
}
