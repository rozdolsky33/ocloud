package resources

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// ListInstances lists all instances in the specified compartment.
func ListInstances(ctx context.Context, provider common.ConfigurationProvider, compartmentID string) error {
	logger.Logger.V(1).Info("ListInstances()")

	// Create a compute client
	_, err := oci.NewComputeClient(provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	fmt.Println("Inside Instance resources running List Instances")
	// Use computeClient to list instances
	// ...

	return nil
}

// FindInstances finds instances by name in the specified compartment.
func FindInstances(ctx context.Context, provider common.ConfigurationProvider, compartmentID, namePattern string, showImageDetails bool) error {
	logger.Logger.V(1).Info("FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails)

	// Create a compute client
	_, err := oci.NewComputeClient(provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	// Create a network client
	_, err = oci.NewNetworkClient(provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	fmt.Println("Finding instances with name pattern:", namePattern)
	// Use computeClient and networkClient to find instances
	// ...

	return nil
}
