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

// CommandRegistry holds information about commands that don't need a full application context
type CommandRegistry struct {
	NoContextCommands map[string]bool
	NoContextFlags    map[string]bool
}

// DefaultRegistry is the global command registry
var DefaultRegistry = &CommandRegistry{
	NoContextCommands: map[string]bool{
		"version": true,
	},
	NoContextFlags: map[string]bool{
		"--version": true,
		"-v":        true,
	},
}

// RegisterNoContextCommand adds a command to the registry of commands that don't need context
func (r *CommandRegistry) RegisterNoContextCommand(cmdName string) {
	r.NoContextCommands[cmdName] = true
}

// RegisterNoContextFlag adds a flag to the registry of flags that don't need context
func (r *CommandRegistry) RegisterNoContextFlag(flagName string) {
	r.NoContextFlags[flagName] = true
}

// IsNoContextCommand checks if the command being run is one that doesn't need context
func (r *CommandRegistry) IsNoContextCommand() bool {
	args := os.Args
	if len(args) < 2 {
		return false
	}

	// Check for direct command match
	if r.NoContextCommands[args[1]] {
		return true
	}

	// Check for flag match
	for _, arg := range args[1:] {
		if r.NoContextFlags[arg] {
			return true
		}
	}

	return false
}

// NewRootCmd creates a new root command with all subcommands attached
func NewRootCmd(appCtx *app.ApplicationContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	// Initialize global flags
	flags.AddGlobalFlags(rootCmd)

	rootCmd.AddCommand(compute.NewComputeCmd(appCtx))

	rootCmd.AddCommand(identity.NewIdentityCmd(appCtx))

	rootCmd.AddCommand(database.NewDatabaseCmd(appCtx))

	rootCmd.AddCommand(network.NewNetworkCmd(appCtx))

	return rootCmd
}

// createRootCmdWithoutContext creates a root command without application context
// This is used for commands that don't need a full context
func createRootCmdWithoutContext() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	// Initialize global flags
	flags.AddGlobalFlags(rootCmd)

	// Add commands that don't need context
	// Currently, only the version command doesn't need context
	rootCmd.AddCommand(version.NewVersionCommand(nil))
	version.AddVersionFlag(rootCmd)

	rootCmd.AddCommand(configuration.NewConfigCmd())

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
	if DefaultRegistry.IsNoContextCommand() {
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
