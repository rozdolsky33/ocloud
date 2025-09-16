package compartment

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	scopeFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/cmd/identity/util"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/compartment"
	"github.com/spf13/cobra"
)

// Long description for the list command
var getLong = `
Get all Compartments in the specified tenancy or compartment with pagination support.

This command displays information about compartments in the current tenancy.
By default, it shows basic compartment information such as name, ID, and description.

The output is paginated, with a default limit of 20 compartments per page. You can navigate
through pages using the --page flag and control the number of compartments per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available compartments in the tenancy
`

// Examples for the get command
var getExamples = `
  # Get all compartments with default pagination (20 per page)
  ocloud identity compartment get

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
			return RunGetCommand(cmd, appCtx)
		},
	}
	// Add flags specific to the get command
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

	scope := util.ResolveScope(cmd)
	parentID := util.ResolveParentID(scope, appCtx)

	logger.LogWithLevel(
		logger.CmdLogger, logger.Debug, "Running compartment get",
		"scope", scope, "parentID", parentID, "json", useJSON,
	)
	return compartment.GetCompartments(appCtx, useJSON, limit, page, parentID)
}
