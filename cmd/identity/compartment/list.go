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

var listLong = `
Interactively browse and search Compartments in the specified tenancy or parent compartment using a TUI.

This command launches a terminal UI that loads available Compartments and lets you:
- Search/filter Compartments as you type
- Navigate the list
- Select a single Compartment to view its details

After you pick a Compartment, the tool prints detailed information about the selected Compartment in a default table view or JSON format if specified with --json.

Scope control:
- Use --scope to choose where to list from: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level listing; it overrides --scope.
- When scope is tenancy, the TUI lists all compartments in the tenancy (including subtree).
- When scope is compartment, the TUI lists only the direct children of the configured compartment.
`

var listExamples = `
  # Launch the interactive compartments browser (default scope: compartment)
  ocloud identity compartment list

  # List at tenancy level (equivalent ways)
  ocloud identity compartment list -T
  ocloud identity compartment list --scope tenancy

  # List only direct children of the configured compartment (explicit)
  ocloud identity compartment list --scope compartment

  # Output selection in JSON format
  ocloud identity compartment list --json
`

// NewListCmd creates a new Cobra command for getting compartments in a specified tenancy or compartment.
// It supports pagination and optional JSON output.
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Compartments in the specified tenancy or compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}

	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)

	return cmd

}

// runListCommand handles the execution of the get command
func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)
	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running compartment list",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return compartment.ListCompartments(appCtx, parentID, useJSON)
}
