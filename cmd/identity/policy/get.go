package policy

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/policy"
	"github.com/spf13/cobra"
)

var getLong = `
Get all Policies in the specified compartment with pagination support.

This command displays information about available Policies in the current compartment.
By default, it shows basic policy information such as name, ID, and description.

The output is paginated, with a default limit of 20 policies per page. You can navigate
through pages using the --page flag and control the number of policies per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available Policies in the compartment
`

var getExamples = `
  # Get all Policies with default pagination (20 per page)
  ocloud identity policy get

  # Get Policies with custom pagination (10 per page, page 2)
  ocloud identity policy get --limit 10 --page 2

  # Get Policies and output in JSON format
  ocloud identity policy get --json

  # Get Policies with custom pagination and JSON output
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

	return cmd

}

// RunGetCommand handles the execution of the get command
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running policy get command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return policy.GetPolicies(appCtx, useJSON, limit, page)
}
