package instance

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociinstance "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
)

// FindInstances finds and displays instances matching a name pattern.
func FindInstances(appCtx *app.ApplicationContext, namePattern string, useJSON, showDetails bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	instanceAdapter := ociinstance.NewAdapter(computeClient, networkClient)
	service := NewService(instanceAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedInstances, err := service.Find(context.Background(), namePattern, showDetails)
	if err != nil {
		return fmt.Errorf("finding instances: %w", err)
	}

	return PrintInstancesInfo(matchedInstances, appCtx, nil, useJSON, showDetails)
}
