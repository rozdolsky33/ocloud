package dynamicgroup

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/dynamicgroup"
)

func SearchDynamicGroups(appCtx *app.ApplicationContext, tenancyOCID string, pattern string, useJSON bool) error {
	ctx := context.Background()
	dgAdapter := dynamicgroup.NewDynamicGroupAdapter(appCtx.IdentityClient, appCtx.Provider)
	service := NewService(dgAdapter, appCtx.Logger, tenancyOCID)

	results, err := service.FuzzySearch(ctx, pattern)
	if err != nil {
		return fmt.Errorf("searching dynamic groups: %w", err)
	}

	return PrintDynamicGroupsInfo(results, appCtx, nil, useJSON)
}
