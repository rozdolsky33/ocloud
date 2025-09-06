package autonomousdb

import (
	databaseFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	"github.com/spf13/cobra"
)

// Long description for the list command
var getLong = `
Fetch Autonomous Databases in the specified compartment with pagination support.

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
var getExamples = `
  # FetchPaginatedClusters all Autonomous Databases with default pagination (20 per page)
  ocloud database autonomous get

  # FetchPaginatedClusters Autonomous Databases with custom pagination (10 per page, page 2)
  ocloud database autonomous get --limit 10 --page 2

  # FetchPaginatedClusters Autonomous Databases and output in JSON format
  ocloud database autonomous get --json

  # FetchPaginatedClusters Autonomous Databases with custom pagination and JSON output
  ocloud database autonomous get --limit 5 --page 3 --json
`

// NewGetCmd creates a "list" subcommand for listing all databases in the specified compartment with pagination support.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all Databases",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, appCtx)
		},
	}

	databaseFlags.LimitFlag.Add(cmd)
	databaseFlags.PageFlag.Add(cmd)
	databaseFlags.AllInfoFlag.Add(cmd)

	return cmd

}

// RunGetCommand handles the execution of the list command
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running autonomous database Get command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, databaseFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, databaseFlags.FlagDefaultPage)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAllInformation, false)
	return autonomousdb.GetAutonomousDatabase(appCtx, useJSON, limit, page, showAll)
}
