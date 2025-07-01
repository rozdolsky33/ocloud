package oke

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func ListClusters(appCtx *app.ApplicationContext, useJSON bool, limit, page int) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing OKE clusters", "limit", limit, "page", page)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating oke cluster service: %w", err)
	}

	ctx := context.Background()
	clusters, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing oke clusters: %w", err)
	}

	err = PrintOKETable(clusters, appCtx, &util.PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	}, useJSON)
	if err != nil {
		return fmt.Errorf("printing clusters: %w", err)
	}

	return nil
}
