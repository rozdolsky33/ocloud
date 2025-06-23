package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func FindCompartments(appCtx *app.ApplicationContext, namePattern string, useJSON bool) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Finding Compartments", "pattern", namePattern)
	return nil
}
