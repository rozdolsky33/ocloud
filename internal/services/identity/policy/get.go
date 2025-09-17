package policy

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/policy"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetPolicies retrieves and displays the policies for a given application context, supporting pagination and JSON output format.
func GetPolicies(appCtx *app.ApplicationContext, useJSON bool, limit, page int, ocid string) error {
	ctx := context.Background()
	policyAdapter := policy.NewAdapter(appCtx.IdentityClient)
	service := NewService(policyAdapter, appCtx.Logger, ocid)
	policies, totalCount, nextPageToken, err := service.FetchPaginatedPolies(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("getting policies: %w", err)
	}
	return PrintPolicyInfo(policies, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
}
