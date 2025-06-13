package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/cmd/compute"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// NewRootCmd creates a new root command with all subcommands attached
func NewRootCmd(appCtx *app.AppContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	// Initialize global flags
	flags.AddGlobalFlags(rootCmd)

	// Add subcommands, passing in the AppContext
	rootCmd.AddCommand(compute.NewComputeCmd(appCtx))

	return rootCmd
}

// Execute runs the root command with the given context
func Execute(ctx context.Context) {
	// Create a temporary root command for bootstrapping
	tempRoot := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	// Initialize global flags on the temporary root
	flags.AddGlobalFlags(tempRoot)

	setLogLevel(tempRoot)

	// One-shot bootstrap of AppContext
	appCtx, err := app.InitApp(ctx, tempRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Create the real root command with the AppContext
	root := NewRootCmd(appCtx)

	// Execute the command
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setLogLevel(tempRoot *cobra.Command) {
	// Parse the flags to get the log level
	tempRoot.ParseFlags(os.Args)

	// Set the logger level from the flag value
	logLevelFlag := tempRoot.PersistentFlags().Lookup(flags.FlagNameLogLevel)
	if logLevelFlag != nil {
		// Use the value from the parsed flag
		logger.LogLevel = logLevelFlag.Value.String()
		if logger.LogLevel == "" {
			// If not set, use the default value
			logger.LogLevel = "info"
		}
	}

	// This is a Hack!
	//Check if --log-level flag is explicitly set in the command line arguments
	// This ensures that the log level is set correctly regardless of whether
	// the full command or shorthand flags are used
	for i, arg := range os.Args {
		if arg == "--log-level" && i+1 < len(os.Args) {
			logger.LogLevel = os.Args[i+1]
			break
		} else if strings.HasPrefix(arg, "--log-level=") {
			logger.LogLevel = strings.TrimPrefix(arg, "--log-level=")
			break
		}
	}
	// Set the colored output from the flag value
	colorFlag := tempRoot.PersistentFlags().Lookup("color")
	if colorFlag != nil {
		// Use the value from the parsed flag
		colorValue := colorFlag.Value.String()
		logger.ColoredOutput = colorValue == "true"
	}

	// This is a Hack!
	// Check if --color flag is explicitly set in the command line arguments
	// This ensures that the color setting is set correctly regardless of whether
	// the full command or shorthand flags are used
	for _, arg := range os.Args {
		if arg == "--color" {
			logger.ColoredOutput = true
			break
		}
	}

	// Initialize logger
	if err := logger.SetLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize package-level logger with the same logger instance
	logger.InitLogger(logger.CmdLogger)
}
