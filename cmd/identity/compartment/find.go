package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/compartment"
	"github.com/spf13/cobra"
)

// NewFindCmd creates a new command for finding image by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find compartment by name pattern",
		Long:          "Find compartment that match the given pattern.",
		Example:       "  ocloud identity compartment find mycompartment\n  ocloud identity compartment find somecompartment \n  ocloud identity compartment find mycomp --json",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running compartment find command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return compartment.FindCompartments(appCtx, namePattern, useJSON)
}
