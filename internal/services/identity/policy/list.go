package policy

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// ListPolicies retrieves and displays the policies for a given application context, supporting pagination and JSON output format.
func ListPolicies(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, logger.Debug, "Listing Policies", "limit", limit, "page", page)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating policy service: %w", err)
	}

	ctx := context.Background()
	policies, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing policies: %w", err)
	}

	// Display policies information with pagination details
	err = PrintPolicyInfo(policies, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)

	if err != nil {
		return fmt.Errorf("printing policies: %w", err)
	}
	logger.Logger.V(logger.Info).Info("Policy list operation completed successfully.")
	return nil
}
