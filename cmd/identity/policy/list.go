package policy

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/policy"
	"github.com/spf13/cobra"
)

// Long description for the list command
var listLong = `
List all Policies in the specified compartment with pagination support.

This command displays information about available Policies in the current compartment.
By default, it shows basic policy information such as name, ID, and description.

The output is paginated, with a default limit of 20 policies per page. You can navigate
through pages using the --page flag and control the number of policies per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available Policies in the compartment
`

// Examples for the list command
var listExamples = `
  # List all Policies with default pagination (20 per page)
  ocloud identity policy list

  # List Policies with custom pagination (10 per page, page 2)
  ocloud identity policy list --limit 10 --page 2

  # List Policies and output in JSON format
  ocloud identity policy list --json

  # List Policies with custom pagination and JSON output
  ocloud identity policy list --limit 5 --page 3 --json
`

// NewListCmd creates a new cobra.Command for listing all policies in a specified tenancy or compartment.
// The command supports pagination through the --limit and --page flags for controlling list size and navigation.
// It also provides optional JSON output for formatted results using the --JSON flag.
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
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running policy list command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return policy.ListPolicies(appCtx, useJSON, limit, page)
}
