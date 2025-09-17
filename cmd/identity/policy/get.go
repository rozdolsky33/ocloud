package policy

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	scopeFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	scopeUtil "github.com/rozdolsky33/ocloud/cmd/shared/scope"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/policy"
	"github.com/spf13/cobra"
)

var getLong = `
Get all Policies in the specified tenancy or compartment with pagination support.

This command displays information about policies either in the current tenancy or within
a specific parent compartment (depending on scope). By default, it shows basic policy
information such as name, ID, and description.

Pagination:
- The output is paginated, with a default limit of 20 policies per page.
- Navigate through pages using the --page flag and control items per page with --limit.

Output formats:
- Use --json (-j) to output the results in JSON format; otherwise a table is shown.

Scope control:
- Use --scope to choose where to list from: "compartment" (default) or "tenancy".
- Use -T/--tenancy-scope as a shortcut to force tenancy-level listing; it overrides --scope.
- When scope is tenancy, the command lists all policies in the tenancy (including subtree, if applicable).
- When scope is compartment, the command lists only the direct children of the configured compartment.`

var getExamples = `
  # Get all policies with default pagination (20 per page)
  ocloud identity policy get

  # List at tenancy level (equivalent ways)
  ocloud identity policy get -T
  ocloud identity policy get --scope tenancy

  # List direct children of the configured compartment (explicit)
  ocloud identity policy get --scope compartment

  # Get policies with custom pagination (10 per page, page 2)
  ocloud identity policy get --limit 10 --page 2

  # Get policies and output in JSON format
  ocloud identity policy get --json

  # Get policies with custom pagination and JSON output
  ocloud identity policy get --limit 5 --page 3 --json
`

// NewGetCmd creates a new cobra.Command for get all policies in a specified tenancy or compartment.
// The command supports pagination through the --limit and --page flags for controlling get size and navigation.
// It also provides optional JSON output for formatted results using the --JSON flag.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Aliases:       []string{"l"},
		Short:         "Get all Policies in the specified tenancy or compartment",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, appCtx)
		},
	}
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)
	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)

	return cmd

}

// RunGetCommand handles the execution of the get command
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running policy get",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return policy.GetPolicies(appCtx, useJSON, limit, page, parentID)
}
