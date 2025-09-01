package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/compartment"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
Find Compartments in the specified tenancy that match the given pattern.

This command searches for compartments whose names match the specified pattern.
By default, it shows basic compartment information such as name, ID, and description
for all matching compartments.

The search is performed using fuzzy matching, which means it will find compartments
even if the pattern is only partially matched. The search is case-insensitive.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command searches across all available compartments in the tenancy
`

// Examples for the find command
var findExamples = `
  # Find compartments with names containing "prod"
  ocloud identity compartment find prod

  # Find compartments with names containing "dev" and output in JSON format
  ocloud identity compartment find dev --json

  # Find compartments with names containing "test" (case-insensitive)
  ocloud identity compartment find test
`

// NewFindCmd creates a new command for finding compartments by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find compartment by name pattern",
		Long:          findLong,
		Example:       findExamples,
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
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running find command", "pattern", namePattern, "json", useJSON)
	return compartment.FindCompartments(appCtx, namePattern, useJSON)
}
