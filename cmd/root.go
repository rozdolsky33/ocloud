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
	"strings"
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

// Execute runs the root command with the given context
func Execute(ctx context.Context) {
	// Create a temporary root command for bootstrapping
	tempRoot := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	flags.AddGlobalFlags(tempRoot)

	setLogLevel(tempRoot)

	appCtx := hasHelpFlag(ctx, tempRoot)

	// Create the real root command with the AppContext
	root := NewRootCmd(appCtx)

	// Execute the command
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func hasHelpFlag(ctx context.Context, tempRoot *cobra.Command) *app.AppContext {
	// Check if help flag is present
	isHelpRequested := false
	for _, arg := range os.Args {
		if arg == flags.FlagPrefixShortHelp || arg == flags.FlagPrefixLongHelp || arg == flags.FlagNameHelp {
			isHelpRequested = true
			break
		}
	}

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
			fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
			os.Exit(1)
		}
	}
	return appCtx
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
			logger.LogLevel = flags.FlagValueInfo
		}
	}

	// This is a Hack!
	//Check if --log-level flag is explicitly set in the command line arguments
	// This ensures that the log level is set correctly regardless of whether
	// the full command or shorthand flags are used
	for i, arg := range os.Args {
		if arg == flags.FlagPrefixLogLevel && i+1 < len(os.Args) {
			logger.LogLevel = os.Args[i+1]
			break
		} else if strings.HasPrefix(arg, flags.FlagPrefixLogLevelEq) {
			logger.LogLevel = strings.TrimPrefix(arg, flags.FlagPrefixLogLevelEq)
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
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize package-level logger with the same logger instance
	logger.InitLogger(logger.CmdLogger)
}
