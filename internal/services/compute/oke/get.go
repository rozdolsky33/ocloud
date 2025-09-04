package oke

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocioke "github.com/rozdolsky33/ocloud/internal/oci/compute/oke"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetClusters retrieves and displays a paginated list of OKE clusters.
func GetClusters(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	containerEngineClient, err := oci.NewContainerEngineClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating container engine client: %w", err)
	}

	clusterAdapter := ocioke.NewAdapter(containerEngineClient)
	service := NewService(clusterAdapter, appCtx.Logger, appCtx.CompartmentID)

	clusters, totalCount, nextPageToken, err := service.FetchPaginatedClusters(context.Background(), limit, page)
	if err != nil {
		return fmt.Errorf("listing clusters: %w", err)
	}

	return PrintOKETable(clusters, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
}
