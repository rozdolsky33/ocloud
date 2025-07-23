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
	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
)

// CommandFactory defines an interface for creating commands
//
//	can be created with or without an application context
type CommandFactory interface {
	// CreateCommand creates a command with the given application context
	// For commands that don't need the context, it will be ignored
	CreateCommand(appCtx *app.ApplicationContext) *cobra.Command
}

// ContextAwareCommandFactory is a factory for commands that require an application context
type ContextAwareCommandFactory struct {
	// createFn is a function that creates a command with an application context
	createFn func(appCtx *app.ApplicationContext) *cobra.Command
	// name is the name of the command, used for logging
	name string
}

// CreateCommand creates a command with the given application context
// If appCtx is nil, it will create a wrapper command that initializes the application context on demand
func (f *ContextAwareCommandFactory) CreateCommand(appCtx *app.ApplicationContext) *cobra.Command {
	// If appCtx is not nil, just create the command directly
	if appCtx != nil {
		return f.createFn(appCtx)
	}

	// Create a minimal ApplicationContext for the temporary command
	// This is safer than passing nil, as some commands might not handle nil gracefully
	minimalAppCtx := &app.ApplicationContext{
		Provider: config.LoadOCIConfig(),
	}

	// Create a temporary instance of the real command to get its properties
	// This is used to copy properties like aliases to the wrapper command
	tempCmd := f.createFn(minimalAppCtx)

	// If appCtx is nil, create a wrapper command that will initialize the application context on demand
	// This allows commands that need the application context to be created even if the application context is nil
	// The application context will be initialized when the command is executed
	wrapperCmd := &cobra.Command{
		Use:           f.name,
		Aliases:       tempCmd.Aliases,
		Short:         tempCmd.Short,
		Long:          tempCmd.Long,
		SilenceUsage:  tempCmd.SilenceUsage,
		SilenceErrors: tempCmd.SilenceErrors,
		Example:       tempCmd.Example,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize the application context
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Create a temporary root command for bootstrapping
			tempRoot := &cobra.Command{
				Use:          "ocloud",
				Short:        "Interact with Oracle Cloud Infrastructure",
				Long:         "",
				SilenceUsage: true,
			}

			flags.AddGlobalFlags(tempRoot)

			// Initialize the application context
			appCtx, err := InitializeAppContext(ctx, tempRoot)
			if err != nil {
				return fmt.Errorf("initializing app context: %w", err)
			}

			// Create the real command with the application context
			realCmd := f.createFn(appCtx)

			// Execute the real command
			if realCmd.RunE != nil {
				return realCmd.RunE(realCmd, args)
			} else if realCmd.Run != nil {
				realCmd.Run(realCmd, args)
				return nil
			} else {
				// If neither RunE nor Run is defined, show help
				return realCmd.Help()
			}
		},
	}

	// Add all subcommands from the temporary command to the wrapper command
	for _, subCmd := range tempCmd.Commands() {
		wrapperCmd.AddCommand(subCmd)
	}

	return wrapperCmd
}

// ContextlessCommandFactory is a factory for commands that don't require an application context
type ContextlessCommandFactory struct {
	// createFn is a function that creates a command without an application context
	createFn func() *cobra.Command
}

// CreateCommand creates a command, ignoring the application context
func (f *ContextlessCommandFactory) CreateCommand(appCtx *app.ApplicationContext) *cobra.Command {
	return f.createFn()
}

// NewRootCmd creates a new root command with all subcommands attached.
// It sets up the command structure for the ocloud CLI tool, including compute, identity,
// database, network, and configuration commands. Each subcommand provides specific
// functionality for interacting with different aspects of Oracle Cloud Infrastructure.
func NewRootCmd(appCtx *app.ApplicationContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "ocloud is a command-line tool for interacting with Oracle Cloud Infrastructure (OCI). It provides commands for managing compute instances, identity resources, databases, networking, and configuration settings.",
		Example:      "  ocloud compute instance list\n  ocloud identity compartment list\n  ocloud database autonomousdb list\n  ocloud network subnet list\n  ocloud config info map-file",
		SilenceUsage: true,
	}

	// Initialize global flags
	flags.AddGlobalFlags(rootCmd)

	// Define command factories
	commandFactories := map[string]CommandFactory{
		"version": &ContextlessCommandFactory{
			createFn: func() *cobra.Command {
				// For version command, we can use nil for appCtx since it handles nil gracefully
				return version.NewVersionCommand(nil)
			},
		},
		"compute": &ContextAwareCommandFactory{
			createFn: compute.NewComputeCmd,
			name:     "compute",
		},
		"identity": &ContextAwareCommandFactory{
			createFn: identity.NewIdentityCmd,
			name:     "identity",
		},
		"database": &ContextAwareCommandFactory{
			createFn: database.NewDatabaseCmd,
			name:     "database",
		},
		"network": &ContextAwareCommandFactory{
			createFn: network.NewNetworkCmd,
			name:     "network",
		},
		"config": &ContextlessCommandFactory{
			createFn: configuration.NewConfigCmd,
		},
	}

	// Add commands using factories
	for _, factory := range commandFactories {
		rootCmd.AddCommand(factory.CreateCommand(appCtx))
	}

	// Add version flag
	version.AddVersionFlag(rootCmd)

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

	// Create the root command without initializing the application context
	// The ContextAwareCommandFactory will initialize the application context on demand
	// for commands that need it
	root := NewRootCmd(nil)

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
