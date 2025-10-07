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

// Long description for the get command
var getLong = `
Get all Compartments in the specified tenancy or compartment with pagination support.

This command displays information about compartments in the current tenancy or within
a specific parent compartment (depending on scope). By default, it shows basic
compartment information such as name, ID, and description.

Pagination:
- The output is paginated, with a default limit of 20 compartments per page.
- Navigate through pages using the --page flag and control items per page with --limit.

Output formats:
- Use --json (-j) to output the results in JSON format; otherwise a table is shown.

Scope control:
- Use --scope to choose where to list from: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level listing; it overrides --scope.
- When scope is tenancy, the command lists all compartments in the tenancy (including subtree).
- When scope is compartment, the command lists only the direct children of the configured compartment.`

// Examples for the get command
var getExamples = `
  # Get all compartments with default pagination (20 per page)
  ocloud identity compartment get

  # List at tenancy level (equivalent ways)
  ocloud identity compartment get -T
  ocloud identity compartment get --scope tenancy

  # List direct children of the configured compartment (explicit)
  ocloud identity compartment get --scope compartment

  # Get compartments with custom pagination (10 per page, page 2)
  ocloud identity compartment get --limit 10 --page 2

  # Get compartments and output in JSON format
  ocloud identity compartment get --json

  # Get compartments with custom pagination and JSON output
  ocloud identity compartment get --limit 5 --page 3 --json
`

// NewGetCmd creates a new Cobra command for getting compartments in a specified tenancy or compartment.
// It supports pagination and optional JSON output.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all Compartments in the specified tenancy or compartment",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}
	scopeFlags.LimitFlag.Add(cmd)
	scopeFlags.PageFlag.Add(cmd)
	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)

	return cmd

}

// runGetCommand handles the execution of the get command
func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, scopeFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, scopeFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running compartment get",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return compartment.GetCompartments(appCtx, useJSON, limit, page, parentID)
}
