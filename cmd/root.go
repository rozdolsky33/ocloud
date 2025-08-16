package cmd

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/cmd/shared/cmdcreate"
	"github.com/rozdolsky33/ocloud/cmd/shared/cmdutil"
	"github.com/rozdolsky33/ocloud/cmd/shared/display"
	cmdlogger "github.com/rozdolsky33/ocloud/cmd/shared/logger"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command with all subcommands attached
func NewRootCmd(appCtx *app.ApplicationContext) *cobra.Command {
	return cmdcreate.CreateRootCmd(appCtx)
}

// Execute runs the root command with the given context.
// It now returns an error instead of exiting directly.
func Execute(ctx context.Context) error {
	// Create a temporary root command for bootstrapping
	tempRoot := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	flags.AddGlobalFlags(tempRoot)

	if err := cmdlogger.SetLogLevel(tempRoot); err != nil {
		return fmt.Errorf("setting log level: %w", err)
	}

	if cmdutil.IsNoContextCommand() {
		root := cmdcreate.CreateRootCmdWithoutContext()

		if cmdutil.IsRootCommandWithoutSubcommands() {
			display.PrintOCIConfiguration()
		}

		if err := root.ExecuteContext(ctx); err != nil {
			return fmt.Errorf("failed to execute root command: %w", err)
		}

		return nil
	}

	appCtx, err := InitializeAppContext(ctx, tempRoot)
	if err != nil {
		return fmt.Errorf("initializing app context: %w", err)
	}

	if cmdutil.IsRootCommandWithoutSubcommands() {
		display.PrintOCIConfiguration()
	}

	// Create the real root command with the ApplicationContext
	root := cmdcreate.CreateRootCmd(appCtx)

	// Set the default behavior to show help
	root.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	// Execute the command
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}

	return nil
}
