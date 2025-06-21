package oke

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func ListClusters(appCtx *app.ApplicationContext, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Listing OKE clusters")

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating oke cluster service: %w", err)
	}

	ctx := context.Background()
	clusters, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("listing oke clusters: %w", err)
	}

	err = PrintOKEInfo(clusters, appCtx, useJSON)
	if err != nil {
		return fmt.Errorf("printing clusters: %w", err)
	}

	return nil
}
