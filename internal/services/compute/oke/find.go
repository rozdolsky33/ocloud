package oke

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocioke "github.com/rozdolsky33/ocloud/internal/oci/compute/oke"
)

// FindClusters finds and displays OKE clusters matching a name pattern.
func FindClusters(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	containerEngineClient, err := oci.NewContainerEngineClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating container engine client: %w", err)
	}

	clusterAdapter := ocioke.NewAdapter(containerEngineClient)
	service := NewService(clusterAdapter, appCtx.Logger, appCtx.CompartmentID)

	matchedClusters, err := service.Find(context.Background(), namePattern)
	if err != nil {
		return fmt.Errorf("finding clusters: %w", err)
	}

	return PrintOKEsInfo(matchedClusters, appCtx, nil, useJSON)
}
