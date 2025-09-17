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

var listLong = `
Interactively browse and search Policies in the specified tenancy or parent compartment using a TUI.

This command launches a terminal UI that loads available policies and lets you:
- Search/filter policies as you type
- Navigate the list
- Select a single policy to view its details

After you pick a policy, the tool prints detailed information about the selected policy in a default table view or JSON format if specified with --json.

Scope control:
- Use --scope to choose where to list from: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level listing; it overrides --scope.
- When scope is tenancy, the TUI lists all policies in the tenancy (including subtree).
- When scope is compartment, the TUI lists only the direct children of the configured compartment.
`

var listExamples = `
  # Launch the interactive policies browser (default scope: compartment)
  ocloud identity policy list

  # List at tenancy level (equivalent ways)
  ocloud identity policy list -T
  ocloud identity policy list --scope tenancy

  # List only direct children of the configured compartment (explicit)
  ocloud identity policy list --scope compartment

  # Output selection in JSON format
  ocloud identity policy list --json
`

// NewListCmd returns "policy list".
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Policies in the specified tenancy or compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}
	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)
	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running policy list",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)

	return policy.ListPolicies(appCtx, useJSON, parentID)
}
