package dynamicgroup

import (
	scopeFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	scopeUtil "github.com/rozdolsky33/ocloud/cmd/shared/scope"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/dynamicgroup"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Dynamic Groups in the specified tenancy or parent compartment using a TUI.

Scope control:
- Use --scope to choose where to list from: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level listing; it overrides --scope.
- Most Dynamic Groups are defined at the tenancy level.
`

var listExamples = `
  # Launch the interactive dynamic groups browser (default scope: compartment)
  ocloud identity dynamic-group list

  # List at tenancy level (recommended for dynamic groups)
  ocloud identity dynamic-group list -T
`

// NewListCmd creates a new Cobra command for listing dynamic groups.
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Dynamic Groups in the specified tenancy or compartment",
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

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running dynamic group list",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)

	return dynamicgroup.ListDynamicGroups(appCtx, parentID, useJSON)
}
