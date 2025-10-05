package policy

import (
	scopeFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	scopeUtil "github.com/rozdolsky33/ocloud/cmd/shared/scope"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/policy"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
FuzzySearch Policies in the specified tenancy or parent compartment that match the given pattern.

This command searches for policies whose names or statements match the specified pattern within
the chosen scope. By default, it searches within the configured parent compartment.

Search behavior:
- Uses a fuzzy/contains-style match; partial and case-insensitive matches are supported.
- Searchable fields include Name, Description, and Statements.

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
  # FuzzySearch policies with names containing "admin" (default scope: compartment)
  ocloud identity policy find admin

  # Search at tenancy level (equivalent ways)
  ocloud identity policy find admin -T
  ocloud identity policy find admin --scope tenancy

  # Search only direct children of the configured compartment (explicit)
  ocloud identity policy find admin --scope compartment

  # FuzzySearch policies with names containing "network" and output in JSON format
  ocloud identity policy find network --json
`

// NewFindCmd creates a new command for finding policies by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "FuzzySearch Policies by name pattern",
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
		logger.CmdLogger, logger.Debug, "Running policy find",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running policy find command", "pattern", namePattern, "json", useJSON)
	return policy.FindPolicies(appCtx, namePattern, useJSON, parentID)
}
