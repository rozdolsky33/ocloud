package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/compute/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// NewFindCmd creates a new command for finding instances by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find instances by name pattern",
		Long:          "Find instances in the specified compartment that match the given pattern.",
		Example:       "  ocloud compute instance find myinstance\n  ocloud compute instance find web-server --image-details\n  ocloud compute instance find api --json",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	// Add flags specific to the find command
	instaceFlags.ImageDetailsFlag.Add(cmd)

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	imageDetails := flags.GetBoolFlag(cmd, flags.FlagNameImageDetails, false)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running instance find command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return instance.FindInstances(appCtx, namePattern, imageDetails, useJSON)
}
