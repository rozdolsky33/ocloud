package heatwave

import (
	databaseFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/heatwavedb"
	"github.com/spf13/cobra"
)

var searchLong = `
Fuzzy Search for HeatWave Databases in the specified compartment.

Search across multiple database attributes including name, OCID, version, shape, network,
and tags. The search uses fuzzy matching to find databases even with typos or partial matches.

Searchable fields include:
  - Name, OCID, State, Description
  - MySQL Version, Shape Name, Storage Size
  - Database Mode, Access Mode, Crash Recovery
  - VCN Name/ID, Subnet Name/ID, IP Address, Hostname Label
  - Network Security Group Names/IDs
  - HeatWave Cluster Size
  - Availability Domain, Fault Domain
  - Tags (both keys and values)
`

var searchExamples = `
  # Search by database name
  ocloud database heatwave search prod-db

  # Search by MySQL version
  ocloud database heatwave search 8.4.6

  # Search by shape name
  ocloud database heatwave search MySQL.4

  # Search by network name
  ocloud database heatwave search prod-vcn

  # Search by IP address
  ocloud database heatwave search 10.0.20

  # Search with JSON output
  ocloud database heatwave search prod-db --json

  # Search with detailed output
  ocloud database heatwave search prod-db --all
`

// NewSearchCmd creates a new command for searching HeatWave Databases.
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for HeatWave Databases",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}
	databaseFlags.AllInfoFlag.Add(cmd)
	return cmd
}

// runSearchCommand handles the execution of the search command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running HeatWave database search command", "searchPattern", namePattern, "json", useJSON, "showAll", showAll)
	return heatwavedb.SearchHeatWaveDatabases(appCtx, namePattern, useJSON, showAll)
}
