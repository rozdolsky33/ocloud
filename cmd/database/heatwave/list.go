package heatwave

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/heatwavedb"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search HeatWave Databases in the specified compartment using a TUI.

This command launches terminal UI that loads available HeatWave Databases and lets you:
- Search/filter HeatWave Database as you type
- Navigate the list
- Select a single HeatWave Databases to view its details

After you pick an HeatWave Database, the tool prints detailed information about the selected HeatWave Database default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive images browser
   ocloud database HeatWave list
   ocloud database HeatWave list --json
`

// NewListCmd creates a new command for listing HeatWave Databases
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all HeatWave Databases",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}
	return cmd

}

// runListCommand handles the execution of the list command
func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running HeatWave database list command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	return heatwavedb.ListHeatWaveDatabases(appCtx, useJSON)
}
