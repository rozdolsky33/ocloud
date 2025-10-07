package instance

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociInst "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
)

// SearchInstances queries and retrieves matching instances based on a fuzzy search pattern.
// It uses OCI clients to fetch instance details and prints results based on the provided flags.
func SearchInstances(appCtx *app.ApplicationContext, search string, useJSON, showDetails bool) error {
	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	instanceAdapter := ociInst.NewAdapter(computeClient, networkClient)
	service := NewService(instanceAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedInstances, err := service.FuzzySearch(context.Background(), search)
	if err != nil {
		return fmt.Errorf("finding instances: %w", err)
	}

	err = PrintInstancesInfo(matchedInstances, appCtx, nil, useJSON, showDetails)
	if err != nil {
		return fmt.Errorf("printing instances: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching instances", "search", search, "matched", len(matchedInstances))
	return nil
}
