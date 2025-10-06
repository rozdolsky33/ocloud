package oke

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocioke "github.com/rozdolsky33/ocloud/internal/oci/compute/oke"
)

// SearchOKEClusters searches for OKE clusters matching a search pattern and displays the results in table or JSON format.
// Parameters:
// - appCtx: The application context containing configuration and dependencies.
// - search: The search string used for matching cluster names.
// - useJSON: A flag indicating whether to display the output in JSON format.
// Returns an error if the search or display operation fails.
func SearchOKEClusters(appCtx *app.ApplicationContext, search string, useJSON bool) error {
	containerEngineClient, err := oci.NewContainerEngineClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating container engine client: %w", err)
	}

	clusterAdapter := ocioke.NewAdapter(containerEngineClient)
	service := NewService(clusterAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedClusters, err := service.FuzzySearch(context.Background(), search)
	if err != nil {
		return fmt.Errorf("searching clusters: %w", err)
	}
	err = PrintOKEsInfo(matchedClusters, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing clusters: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching clusters", "search", search, "matched", len(matchedClusters))
	return nil
}
