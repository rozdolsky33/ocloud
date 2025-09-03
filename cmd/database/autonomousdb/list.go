package autonomousdb

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	"github.com/spf13/cobra"
)

// Long description for the list command
var listLong = `
FetchPaginatedInstances all Autonomous Databases in the specified compartment with pagination support.

This command displays information about available Autonomous Databases in the current compartment.
By default, it shows basic database information such as name, ID, state, and workload type.

The output is paginated, with a default limit of 20 databases per page. You can navigate
through pages using the --page flag and control the number of databases per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available Autonomous Databases in the compartment
`

// Examples for the list command
var listExamples = `
  # FetchPaginatedInstances all Autonomous Databases with default pagination (20 per page)
  ocloud database autonomous list

  # FetchPaginatedInstances Autonomous Databases with custom pagination (10 per page, page 2)
  ocloud database autonomous list --limit 10 --page 2

  # FetchPaginatedInstances Autonomous Databases and output in JSON format
  ocloud database autonomous list --json

  # FetchPaginatedInstances Autonomous Databases with custom pagination and JSON output
  ocloud database autonomous list --limit 5 --page 3 --json
`

// NewListCmd creates a "list" subcommand for listing all databases in the specified compartment with pagination support.
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "FetchPaginatedInstances all Databases in the specified compartment",
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
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	return autonomousdb.ListAutonomousDatabase(appCtx, useJSON, limit, page)
}
