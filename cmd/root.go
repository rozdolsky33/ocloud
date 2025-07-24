package cmd

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd/compute"
	"github.com/rozdolsky33/ocloud/cmd/configuration"
	"github.com/rozdolsky33/ocloud/cmd/database"
	"github.com/rozdolsky33/ocloud/cmd/identity"
	"github.com/rozdolsky33/ocloud/cmd/network"
	"github.com/rozdolsky33/ocloud/cmd/version"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
	"os"
)

// noContextCommandChecker provides functionality to check if a command doesn't need a full application context
// This is a simplified version of the previous CommandRegistry, removing unused methods
func isNoContextCommand() bool {
	args := os.Args
	if len(args) < 2 {
		return false
	}

	// Commands that don't need context
	noContextCommands := map[string]bool{
		"version": true,
		"config":  true,
	}

	// Flags that don't need context
	noContextFlags := map[string]bool{
		"--version": true,
		"-v":        true,
	}

	// Check for direct command match
	if noContextCommands[args[1]] {
		return true
	}

	// Check for flag match
	for _, arg := range args[1:] {
		if noContextFlags[arg] {
			return true
		}
	}

	return false
}

// createRootCmd creates a root command with or without application context
// If appCtx is nil, only commands that don't need context are added
// If appCtx is not nil, all commands are added
func createRootCmd(appCtx *app.ApplicationContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	// Initialize global flags
	flags.AddGlobalFlags(rootCmd)

	// Add commands that don't need context
	rootCmd.AddCommand(version.NewVersionCommand())
	version.AddVersionFlag(rootCmd)
	rootCmd.AddCommand(configuration.NewConfigCmd())

	// If appCtx is not nil, add commands that need context
	if appCtx != nil {
		rootCmd.AddCommand(compute.NewComputeCmd(appCtx))
		rootCmd.AddCommand(identity.NewIdentityCmd(appCtx))
		rootCmd.AddCommand(database.NewDatabaseCmd(appCtx))
		rootCmd.AddCommand(network.NewNetworkCmd(appCtx))
	}

	return rootCmd
}

// NewRootCmd creates a new root command with all subcommands attached
func NewRootCmd(appCtx *app.ApplicationContext) *cobra.Command {
	rootCmd := createRootCmd(appCtx)
	return rootCmd
}

// createRootCmdWithoutContext creates a root command without application context
// This is used for commands that don't need a full context
func createRootCmdWithoutContext() *cobra.Command {
	rootCmd := createRootCmd(nil)

	// Set the default behavior to show help
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	return rootCmd
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

	if err := setLogLevel(tempRoot); err != nil {
		return fmt.Errorf("setting log level: %w", err)
	}

	// Check if we're running a command that doesn't need context
	if isNoContextCommand() {
		// Create a root command without application context
		root := createRootCmdWithoutContext()

		// Execute the command
		if err := root.ExecuteContext(ctx); err != nil {
			return fmt.Errorf("failed to execute root command: %w", err)
		}

		return nil
	}

	// For all other commands, initialize the application context
	appCtx, err := InitializeAppContext(ctx, tempRoot)
	if err != nil {
		return fmt.Errorf("initializing app context: %w", err)
	}

	// Create the real root command with the ApplicationContext
	root := NewRootCmd(appCtx)

	// Switch to RunE for the root command
	root.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help() // The default behavior is to show help
	}

	// Execute the command
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}

	return nil
}
