package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// Dedicated documentation for the list command (separate from get)
var listLong = `
Interactively browse and search instances in the specified compartment using a TUI.

This command launches terminal UI that loads available instances and lets you:
- Search/filter instance as you type
- Navigate the list
- Select a single instance to view its details

After you pick an instance, the tool prints detailed information about the selected instance default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive instance browser
  ocloud compute instance list
  ocloud compute instance list --json
`

// NewListCmd creates a new command for listing instances
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Instances",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	instaceFlags.AllInfoFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running instance list command in", "compartment", appCtx.CompartmentName, useJSON)
	return instance.ListInstances(appCtx, useJSON)
}
