package cmd

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
	"os"
)

// InitializeAppContext checks for help-related flags and initializes the AppContext accordingly.
// It returns an error instead of exiting directly.
func InitializeAppContext(ctx context.Context, tempRoot *cobra.Command) (*app.AppContext, error) {
	// Check if a help flag is present
	isHelpRequested := HasHelpFlag(os.Args)

	var appCtx *app.AppContext
	var err error

	if isHelpRequested {
		// If help is requested, create a minimal AppContext without cloud configuration
		appCtx = &app.AppContext{
			Logger:          logger.CmdLogger,
			CompartmentName: flags.FlagValueHelpMode, // Set a dummy value to avoid nil pointer issues
		}
	} else {
		// One-shot bootstrap of AppContext
		appCtx, err = app.InitApp(ctx, tempRoot)
		if err != nil {
			return nil, fmt.Errorf("initializing application: %w", err)
		}
	}
	return appCtx, nil
}

// HasHelpFlag checks if any help-related flags are present in the arguments.
func HasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == flags.FlagPrefixShortHelp || arg == flags.FlagPrefixLongHelp || arg == flags.FlagNameHelp {
			return true
		}
	}
	return false
}
