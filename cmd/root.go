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
	// If no arguments are provided (just the program name), we don't need context
	// This avoids initialization when just displaying help/usage information
	if len(args) < 2 {
		return true
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

	if noContextCommands[args[1]] {
		return true
	}

	for _, arg := range args[1:] {
		if noContextFlags[arg] {
			return true
		}
	}

	return false
}

// isRootCommandWithoutSubcommands checks if the command being executed is the root command without any subcommands or flags
// This is used to determine whether to display the banner and configuration details
func isRootCommandWithoutSubcommands() bool {
	args := os.Args

	// Only display the banner and configuration details when running the root command without any subcommands or flags
	// This means only when running "./ocloud" with no additional arguments
	return len(args) == 1
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

	// Add placeholder commands for help display
	addPlaceholderCommands(rootCmd)

	rootCmd.SetHelpTemplate(`
{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	// Set the default behavior to show help
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	return rootCmd
}

// addPlaceholderCommands adds placeholder commands that will be displayed in help
// but will show a message about needing to initialize if they're actually run
func addPlaceholderCommands(rootCmd *cobra.Command) {
	// Add compute command
	computeCmd := &cobra.Command{
		Use:   "compute",
		Short: "Manage OCI compute services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this command requires application initialization")
		},
	}
	rootCmd.AddCommand(computeCmd)

	// Add identity command
	identityCmd := &cobra.Command{
		Use:   "identity",
		Short: "Manage OCI identity services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this command requires application initialization")
		},
	}
	rootCmd.AddCommand(identityCmd)

	// Add database command
	databaseCmd := &cobra.Command{
		Use:   "database",
		Short: "Manage OCI Database services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this command requires application initialization")
		},
	}
	rootCmd.AddCommand(databaseCmd)

	// Add network command
	networkCmd := &cobra.Command{
		Use:   "network",
		Short: "Manage OCI networking services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this command requires application initialization")
		},
	}
	rootCmd.AddCommand(networkCmd)
}

// displayConfigurationDetails displays the current configuration details
// and checks if required environment variables are set
func displayConfigurationDetails() {
	// Display the banner in plain color
	fmt.Println(" ██████╗  ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ")
	fmt.Println("██╔═══██╗██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗")
	fmt.Println("██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║")
	fmt.Println("██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║")
	fmt.Println("╚██████╔╝╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝")
	fmt.Println(" ╚═════╝  ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝")
	fmt.Println()

	fmt.Println("\033[1mConfiguration Details:\033[0m")

	// Check OCI_CLI_PROFILE
	profile := os.Getenv("OCI_CLI_PROFILE")
	if profile == "" {
		fmt.Println("  \033[33mOCI_CLI_PROFILE\033[0m: \033[31mNot set - Please set profile\033[0m")
	} else {
		fmt.Printf("  \033[33mOCI_CLI_PROFILE\033[0m: %s\n", profile)
	}

	// Check OCI_TENANCY_NAME
	tenancyName := os.Getenv(flags.EnvOCITenancyName)
	if tenancyName == "" {
		fmt.Println("  \033[33mOCI_TENANCY_NAME\033[0m: \033[31mNot set - Please set tenancy\033[0m")
	} else {
		fmt.Printf("  \033[33mOCI_TENANCY_NAME\033[0m: %s\n", tenancyName)
	}

	// Check OCI_COMPARTMENT
	compartment := os.Getenv(flags.EnvOCICompartment)
	if compartment == "" {
		fmt.Println("  \033[33mOCI_COMPARTMENT\033[0m: \033[31mNot set - Please set compartmen name\033[0m")
	} else {
		fmt.Printf("  \033[33mOCI_COMPARTMENT\033[0m: %s\n", compartment)
	}

	fmt.Println()
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

		// Display configuration details only for the root command without subcommands
		if isRootCommandWithoutSubcommands() {
			displayConfigurationDetails()
		}

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

	// Display configuration details only for the root command without subcommands
	if isRootCommandWithoutSubcommands() {
		displayConfigurationDetails()
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
