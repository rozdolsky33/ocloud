package compartment

import (
	scopeFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	scopeUtil "github.com/rozdolsky33/ocloud/cmd/shared/scope"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/compartment"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
FuzzySearch Compartments in the specified tenancy or parent compartment that match the given pattern.

This command searches for compartments whose names match the specified pattern. By default, it
shows basic compartment information such as name, ID, and description for all matching compartments
within the chosen scope.

The search is performed using fuzzy matching, which means it will find compartments even if the
pattern is only partially matched. The search is case-insensitive.

Scope control:
- Use --scope to choose where to search: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level search; it overrides --scope.
- When scope is tenancy, the command searches across all compartments in the tenancy (including subtree).
- When scope is compartment, the command searches only the direct children of the configured compartment.

Additional Information:
- Use --json (-j) to output the results in JSON format
`

// Examples for the find command
var findExamples = `
  # FuzzySearch compartments with names containing "prod" (default scope: compartment)
  ocloud identity compartment find prod

  # Search at tenancy level (equivalent ways)
  ocloud identity compartment find prod -T
  ocloud identity compartment find prod --scope tenancy

  # Search only direct children of the configured compartment (explicit)
  ocloud identity compartment find prod --scope compartment

  # FuzzySearch compartments with names containing "dev" and output in JSON format
  ocloud identity compartment find dev --json
`

// NewFindCmd creates a new command for finding compartments by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "FuzzySearch compartment by name pattern",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)
	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running compartment find",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return compartment.FindCompartments(appCtx, namePattern, useJSON, parentID)
}
