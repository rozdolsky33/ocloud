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

var getLong = `
Get all Dynamic Groups in the specified tenancy or compartment with pagination support.

Pagination:
- The output is paginated, with a default limit of 20 per page.
- Navigate through pages using the --page flag and control items per page with --limit.
`

var getExamples = `
  # Get all dynamic groups with default pagination
  ocloud identity dynamic-group get

  # List at tenancy level
  ocloud identity dynamic-group get -T

  # Get a specific dynamic group by OCID
  ocloud identity dynamic-group get ocid1.dynamicgroup.oc1..example
`

// NewGetCmd creates a new Cobra command for getting dynamic groups.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all Dynamic Groups in the specified tenancy or compartment",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx, args)
		},
	}
	scopeFlags.LimitFlag.Add(cmd)
	scopeFlags.PageFlag.Add(cmd)
	scopeFlags.ScopeFlag.Add(cmd)
	scopeFlags.TenancyScopeFlag.Add(cmd)

	return cmd
}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext, args []string) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	if len(args) > 0 {
		return dynamicgroup.GetDynamicGroup(appCtx, args[0], useJSON)
	}

	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, scopeFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, scopeFlags.FlagDefaultPage)

	scope := scopeUtil.ResolveScope(cmd)
	parentID := scopeUtil.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running dynamic group get",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return dynamicgroup.GetDynamicGroups(appCtx, useJSON, limit, page, parentID)
}
