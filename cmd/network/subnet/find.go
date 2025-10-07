package subnet

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/network/subnet"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
FuzzySearch Subnets in the specified tenancy or compartment that match the given pattern.

This command searches for subnets whose names match the specified pattern.
By default, it shows detailed subnet information such as name, ID, CIDR block,
and whether public IP addresses are allowed for all matching subnets.

The search is performed using fuzzy matching, which means it will find subnets
even if the pattern is only partially matched. The search is case-insensitive.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command searches across all available subnets in the compartment
`

// Examples for the find command
var findExamples = `
  # FuzzySearch subnets with names containing "prod"
  ocloud network subnet find prod

  # FuzzySearch subnets with names containing "dev" and output in JSON format
  ocloud network subnet find dev --json

  # FuzzySearch subnets with names containing "test" (case-insensitive)
  ocloud network subnet find test
`

// NewFindCmd creates a new command for finding subnets by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "FuzzySearch Subnets by name pattern",
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
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running subnet find command", "pattern", namePattern, "json", useJSON)
	return subnet.FindSubnets(appCtx, namePattern, useJSON)
}
