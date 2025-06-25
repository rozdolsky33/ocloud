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
List all Autonomous Databases in the specified compartment with pagination support.

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
  # List all Autonomous Databases with default pagination (20 per page)
  ocloud database autonomousdb list

  # List Autonomous Databases with custom pagination (10 per page, page 2)
  ocloud database autonomousdb list --limit 10 --page 2

  # List Autonomous Databases and output in JSON format
  ocloud database autonomousdb list --json

  # List Autonomous Databases with custom pagination and JSON output
  ocloud database autonomousdb list --limit 5 --page 3 --json
`

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Databases in the specified compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}
	// Add flags specific to the list command
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	// Get pagination parameters
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	return autonomousdb.ListAutonomousDatabase(appCtx, useJSON, limit, page)
}
