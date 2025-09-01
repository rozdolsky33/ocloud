package instance

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociInst "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListInstances retrieves and displays a paginated list of instances.
func ListInstances(appCtx *app.ApplicationContext, useJSON bool, limit, page int, showDetails bool) error {
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

	instances, totalCount, nextPageToken, err := service.List(context.Background(), limit, page)
	if err != nil {
		return fmt.Errorf("listing instances: %w", err)
	}

	return PrintInstancesInfo(instances, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON, showDetails)
}
