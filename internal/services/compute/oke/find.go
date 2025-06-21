package oke

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func FindClusters(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding OKE clusters", "pattern", namePattern)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating oke cluster service: %w", err)
	}

	ctx := context.Background()
	clusters, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding oke clusters: %w", err)
	}

	// Display matched clusters
	if len(clusters) == 0 {
		if useJSON {
			// Return an empty JSON array if no clusters found
			fmt.Fprintln(appCtx.Stdout, `{"clusters": []}`)
		} else {
			fmt.Fprintf(appCtx.Stdout, "No clusters found matching pattern: %s\n", namePattern)
		}
		return nil
	}

	err = PrintOKEInfo(clusters, appCtx, useJSON)
	if err != nil {
		return fmt.Errorf("printing clusters: %w", err)
	}

	return nil
}
