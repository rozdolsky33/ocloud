package oke

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociOke "github.com/rozdolsky33/ocloud/internal/oci/compute/oke"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

// ListClusters lists all OKE clusters in the tenancy.
func ListClusters(appCtx *app.ApplicationContext, useJSON bool) error {
	ctx := context.Background()
	containerEngineClient, err := oci.NewContainerEngineClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating container engine client: %w", err)
	}

	clusterAdapter := ociOke.NewAdapter(containerEngineClient)
	service := NewService(clusterAdapter, appCtx.Logger, appCtx.CompartmentID)

	clusters, err := service.ListClusters(ctx)
	if err != nil {
		return fmt.Errorf("listing allClusters: %w", err)
	}

	// TUI
	model := ociOke.NewImageListModel(clusters)
	id, err := listx.Run(model)
	if err != nil {
		if errors.Is(err, listx.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting image: %w", err)
	}

	cluster, err := service.clusterRepo.GetCluster(ctx, id)
	if err != nil {
		return fmt.Errorf("getting image: %w", err)
	}

	return PrintOKEInfo(appCtx, cluster, useJSON)

}
