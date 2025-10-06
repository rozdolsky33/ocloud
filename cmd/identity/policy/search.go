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

// Long description for the search command
var findLong = `
Fuzzy search Policies in the specified tenancy or parent compartment that match the given pattern.

This command searches for policies whose fields match the specified pattern within
the chosen scope. By default, it searches within the configured parent compartment.

Search behavior:
- Uses a combination of fuzzy, prefix, and substring matching; partial and case-insensitive matches are supported.
- You can search using any of the following fields (partial matches are supported):

Searchable fields:
- Name: Policy name
- Description: Policy description
- OCID: Policy OCID
- Statements: Policy statements (combined)
- TagsKV: Tags in key:value or namespace.key:value format
- TagsVal: Tag values only (without keys)

Scope control:
- Use --scope to choose where to search: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level search; it overrides --scope.
- When scope is tenancy, the command searches across all compartments in the tenancy (including subtree).
- When scope is compartment, the command searches only the direct children of the configured compartment.

Additional Information:
- Use --json (-j) to output the results in JSON format
`

// Examples for the search command
var findExamples = `
  # FuzzySearch policies with names containing "admin" (default scope: compartment)
  ocloud identity policy search admin

  # Search at tenancy level (equivalent ways)
  ocloud identity policy search admin -T
  ocloud identity policy search admin --scope tenancy

  # Search only direct children of the configured compartment (explicit)
  ocloud identity policy search admin --scope compartment

  # FuzzySearch policies with names containing "network" and output in JSON format
  ocloud identity policy search network --json
`

// NewSearchCmd creates a new command for finding policies by name pattern
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for Policies",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}

	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)

	return cmd
}

// RunFindCommand handles the execution of the find command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	search := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running policy search",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running policy search command", "search", search, "json", useJSON)
	return policy.SearchPolicies(appCtx, search, useJSON, parentID)
}
