package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Autonomous Databases in the specified compartment using a TUI.

This command launches terminal UI that loads available Autonomous Databases and lets you:
- Search/filter Autonomous Database as you type
- Navigate the list
- Select a single Autonomous Databases to view its details

After you pick an Autonomous Database, the tool prints detailed information about the selected Autonomous Database default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive images browser
   ocloud database autonomous list
   ocloud database autonomous list --json
`

// NewListCmd creates a new command for listing Autonomous Databases
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Autonomous Databases",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}
	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running autonomous database list command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	return autonomousdb.ListAutonomousDatabases(appCtx, useJSON)
}
