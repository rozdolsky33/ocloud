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

// Long description for the search command
var searchLong = `
Search for compartments in the specified scope that match the given pattern.

The search uses a fuzzy, prefix, and substring matching algorithm across multiple indexed fields.
You can search using any of the following fields (partial matches are supported):

Searchable fields:
- Name: Compartment display name
- Description: Compartment description
- OCID: Compartment OCID
- State: Lifecycle state (e.g., ACTIVE)
- TagsKV: All tags in key=value form, flattened
- TagsVal: Only tag values (e.g., "prod")

The search pattern is case-insensitive. For very specific inputs (like full OCID),
the search first tries exact and substring matches; otherwise it falls back to broader fuzzy search.

Scope control:
- Use --scope to choose where to search: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope to force tenancy-level search (overrides --scope).
- When scope is tenancy, the command searches across all compartments in the tenancy (including subtree).
- When scope is compartment, the command searches only the direct children of the configured compartment.

Output control:
- Use --json (-j) to output the results in JSON format
`

// Examples for the search command
var searchExamples = `
  # Search by display name (substring)
  ocloud identity compartment search prod

  # Search across the entire tenancy
  ocloud identity compartment search prod -T
  ocloud identity compartment search prod --scope tenancy

  # Search only direct children of the configured compartment
  ocloud identity compartment search dev --scope compartment

  # Search by OCID (exact)
  ocloud identity compartment search ocid1.compartment.oc1..aaaa...

  # Search by tag value only (TagsVal)
  ocloud identity compartment search prod

  # Output in JSON format
  ocloud identity compartment search finance --json
`

// NewFindCmd creates a new command for finding compartments by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for Compartments",
		Long:          searchLong,
		Example:       searchExamples,
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
		logger.CmdLogger, logger.Debug, "Running compartment search",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return compartment.SearchCompartments(appCtx, namePattern, useJSON, parentID)
}
