package subnet

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindSubnets finds and displays subnets based on the provided name pattern, optionally in JSON format.
// It uses the application context for accessing configurations, logging, and output streams.
// Returns an error if any operation fails during subnet retrieval or output generation.
func FindSubnets(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Subnets", "pattern", namePattern)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating subnet service: %w", err)
	}

	ctx := context.Background()
	matchedSubnets, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched subnets: %w", err)
	}

	err = PrintSubnetInfo(matchedSubnets, appCtx, useJSON)
	if err != nil {
		return fmt.Errorf("printing matched subnets: %w", err)
	}

	return nil
}
