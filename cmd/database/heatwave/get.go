package heatwave

import (
	databaseFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/heatwavedb"
	"github.com/spf13/cobra"
)

// Long description for the list command
var getLong = `
Fetch HeatWave Databases in the specified compartment with pagination support.

This command displays information about available HeatWave Databases in the current compartment.
By default, it shows basic database information such as name, ID, state, and workload type.

The output is paginated, with a default limit of 20 databases per page. You can navigate
through pages using the --page flag and control the number of databases per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available HeatWave Databases in the compartment
`

// Examples for the list command
var getExamples = `
  # Get all HeatWave Databases with default pagination (20 per page)
  ocloud database heatwave get

  # Get HeatWave Databases with custom pagination (10 per page, page 2)
  ocloud database heatwave get --limit 10 --page 2

  # Get HeatWave Databases and output in JSON format
  ocloud database heatwave get --json

  # Get HeatWave Databases with custom pagination and JSON output
  ocloud database heatwave get --limit 5 --page 3 --json
`

// NewGetCmd creates a "list" subcommand for listing all databases in the specified compartment with pagination support.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all HeatWave Databases",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}

	databaseFlags.LimitFlag.Add(cmd)
	databaseFlags.PageFlag.Add(cmd)
	databaseFlags.AllInfoFlag.Add(cmd)

	return cmd

}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running HeatWave database Get command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, databaseFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, databaseFlags.FlagDefaultPage)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	return heatwavedb.GetHeatWaveDatabase(appCtx, useJSON, limit, page, showAll)
}
