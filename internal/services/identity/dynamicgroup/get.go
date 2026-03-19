package dynamicgroup

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/dynamicgroup"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetDynamicGroups retrieves and displays a paginated list of dynamic groups.
func GetDynamicGroups(appCtx *app.ApplicationContext, useJSON bool, limit, page int, tenancyOCID string) error {
	ctx := context.Background()
	dgAdapter := dynamicgroup.NewDynamicGroupAdapter(appCtx.IdentityClient, appCtx.Provider)
	service := NewService(dgAdapter, appCtx.Logger, tenancyOCID)

	dynamicGroups, totalCount, nextPageToken, err := service.FetchPaginateDynamicGroups(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing dynamic groups: %w", err)
	}

	return PrintDynamicGroupsTable(dynamicGroups, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
}

// GetDynamicGroup retrieves and displays a single dynamic group by its OCID.
func GetDynamicGroup(appCtx *app.ApplicationContext, ocid string, useJSON bool) error {
	ctx := context.Background()
	dgAdapter := dynamicgroup.NewDynamicGroupAdapter(appCtx.IdentityClient, appCtx.Provider)
	service := NewService(dgAdapter, appCtx.Logger, "")

	dg, err := service.dynamicGroupRepo.GetDynamicGroup(ctx, ocid)
	if err != nil {
		return fmt.Errorf("getting dynamic group: %w", err)
	}

	return PrintDynamicGroupInfo(dg, appCtx, useJSON)
}
