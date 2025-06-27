package policy

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// FindPolicies retrieves and processes policies matching the provided name pattern within the application context.
// appCtx represents the application context with the necessary clients and configurations.
// namePattern specifies the pattern to filter policy names.
// useJSON determines whether the output should be formatted as JSON.
// Returns an error if policy retrieval or processing fails.
func FindPolicies(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Policies", "pattern", namePattern)

	service, err := NewService(appCtx)
	if err != nil {
		return fmt.Errorf("creating policies service: %w", err)
	}

	ctx := context.Background()
	matchedPolicies, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched policies: %w", err)
	}

	err = PrintPolicyInfo(matchedPolicies, appCtx, nil, useJSON)
	if err != nil {
		return fmt.Errorf("printing matched policies: %w", err)
	}

	return nil
}
