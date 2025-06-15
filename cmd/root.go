package cmd

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd/compute"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
	"os"
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

	// Add a custom help flag with a more descriptive message
	rootCmd.Flags().BoolP(flags.FlagNameHelp, flags.FlagShortHelp, false, flags.FlagDescHelp)
	_ = rootCmd.Flags().SetAnnotation(flags.FlagNameHelp, flags.CobraAnnotationKey, []string{flags.FlagValueTrue})

	// Add subcommands, passing in the AppContext
	rootCmd.AddCommand(compute.NewComputeCmd(appCtx))

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

	appCtx, err := initializeAppContext(ctx, tempRoot)
	if err != nil {
		return fmt.Errorf("initializing app context: %w", err)
	}

	// Create the real root command with the AppContext
	root := NewRootCmd(appCtx)

	// Add PersistentPreRunE to handle setup before any command
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Optional: more setup before any command
		return nil
	}

	// Switch to RunE for the root command
	root.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help() // Default behavior is to show help
	}

	// Execute the command
	if err := root.ExecuteContext(ctx); err != nil {
		return err // Will be handled in main
	}

	return nil
}

// initializeAppContext checks for help-related flags and initializes the AppContext accordingly.
// It returns an error instead of exiting directly.
func initializeAppContext(ctx context.Context, tempRoot *cobra.Command) (*app.AppContext, error) {
	// Check if a help flag is present
	isHelpRequested := hasHelpFlag(os.Args)

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

// hasHelpFlag checks if any help-related flags are present in the arguments.
func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == flags.FlagPrefixShortHelp || arg == flags.FlagPrefixLongHelp || arg == flags.FlagNameHelp {
			return true
		}
	}
	return false
}

// // Hacky!
// setLogLevel sets the logging level and colored output based on command-line flags or default values.
// It ensures consistent log settings, initializes the logger, and applies settings globally.
func setLogLevel(tempRoot *cobra.Command) error {
	// Parse the flags to get the log level
	tempRoot.ParseFlags(os.Args)

	// Check for debug flag first - it takes precedence over log-level
	debugFlag := tempRoot.PersistentFlags().Lookup(flags.FlagNameDebug)
	if debugFlag != nil && debugFlag.Value.String() == flags.FlagValueTrue {
		// If debug flag is set, set log level to debug
		logger.LogLevel = "debug"
	} else {
		// Otherwise, use the log-level flag
		logLevelFlag := tempRoot.PersistentFlags().Lookup(flags.FlagNameLogLevel)
		if logLevelFlag != nil {
			// Use the value from the parsed flag
			logger.LogLevel = logLevelFlag.Value.String()
			if logger.LogLevel == "" {
				// If not set, use the default value
				logger.LogLevel = flags.FlagValueInfo
			}
		}

		// This is a Hack!
		// Check if --log-level flag is explicitly set in the command line arguments
		// This ensures that the log level is set correctly regardless of whether
		// the full command or shorthand flags are used
		//for i, arg := range os.Args {
		//	if arg == flags.FlagPrefixLogLevel && i+1 < len(os.Args) {
		//		logger.LogLevel = os.Args[i+1]
		//		break
		//	} else if strings.HasPrefix(arg, flags.FlagPrefixLogLevelEq) {
		//		logger.LogLevel = strings.TrimPrefix(arg, flags.FlagPrefixLogLevelEq)
		//		break
		//	}
		//}
	}

	// This is a Hack!
	// Check if -d or --debug flag is explicitly set in the command line arguments
	// This ensures that debug mode is set correctly regardless of whether
	// the full command or shorthand flags are used
	for _, arg := range os.Args {
		if arg == flags.FlagPrefixDebug || arg == flags.FlagPrefixShortDebug {
			logger.LogLevel = "debug"
			break
		}
	}

	// Set the colored output from the flag value
	colorFlag := tempRoot.PersistentFlags().Lookup(flags.FlagNameColor)
	if colorFlag != nil {
		// Use the value from the parsed flag
		colorValue := colorFlag.Value.String()
		logger.ColoredOutput = colorValue == flags.FlagValueTrue
	}

	// This is a Hack!
	// Check if --color flag is explicitly set in the command line arguments
	// This ensures that the color setting is set correctly regardless of whether
	// the full command or shorthand flags are used
	for _, arg := range os.Args {
		if arg == flags.FlagPrefixColor {
			logger.ColoredOutput = true
			break
		}
	}

	// Initialize logger
	if err := logger.SetLogger(); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	// Initialize package-level logger with the same logger instance
	logger.InitLogger(logger.CmdLogger)

	return nil
}
