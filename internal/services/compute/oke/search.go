package oke

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocioke "github.com/rozdolsky33/ocloud/internal/oci/compute/oke"
)

// SearchOKEClusters finds and displays OKE clusters matching a name pattern.
func SearchOKEClusters(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	containerEngineClient, err := oci.NewContainerEngineClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating container engine client: %w", err)
	}

	clusterAdapter := ocioke.NewAdapter(containerEngineClient)
	service := NewService(clusterAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedClusters, err := service.FuzzySearch(context.Background(), namePattern)
	if err != nil {
		return fmt.Errorf("finding clusters: %w", err)
	}
	err = PrintOKEsInfo(matchedClusters, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing clusters: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching clusters", "searchPattern", namePattern, "matched", len(matchedClusters))
	return nil
}
