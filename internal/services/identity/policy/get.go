package policy

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetPolicies retrieves and displays the policies for a given application context, supporting pagination and JSON output format.
func GetPolicies(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Getting Policies", "limit", limit, "page", page)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating policy service: %w", err)
	}

	ctx := context.Background()
	policies, totalCount, nextPageToken, err := service.List(ctx, limit, page)
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
