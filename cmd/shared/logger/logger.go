package logger

import (
	"fmt"
	"os"

	"github.com/rozdolsky33/ocloud/cmd/version"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
)

// SetLogLevel sets the logging level and colored output based on command-line flags or default values.
func SetLogLevel(tempRoot *cobra.Command) error {
	for _, arg := range os.Args {
		if arg == flags.FlagPrefixVersion || arg == flags.FlagPrefixShortVersion {
			version.PrintVersion()
			os.Exit(0)
		}
	}
	tempRoot.ParseFlags(os.Args)
	// Parse the flags to get the log level Should be approach, but for some reason it prevents parsing flags and give an error
	//if err: = tempRoot.ParseFlags(os.Args); err != nil {
	//	return fmt.Errorf("parsing flags: %w", err)
	//}

	// Check for a debug flag first - it takes precedence over log-level
	debugFlag := tempRoot.PersistentFlags().Lookup(flags.FlagNameDebug)
	if debugFlag != nil && debugFlag.Value.String() == flags.FlagValueTrue {
		// If a debug flag is set, set the log level to debug
		logger.LogLevel = flags.FlagNameDebug
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
	}

	// This is a Hack!
	// Check if -d or --debug flag is explicitly set in the command line arguments
	// This ensures that debug mode is set correctly regardless of whether
	// the full command or shorthand flags are used
	for _, arg := range os.Args {
		if arg == flags.FlagPrefixDebug || arg == flags.FlagPrefixShortDebug {
			logger.LogLevel = flags.FlagNameDebug
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

	logger.InitLogger(logger.CmdLogger)

	return nil
}
