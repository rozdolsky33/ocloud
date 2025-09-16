package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/compartment"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Compartments in the specified tenancy or parent compartment using a TUI.

This command launches terminal UI that loads available Compartments and lets you:
- Search/filter Compartment as you type
- Navigate the list
- Select a single Compartment to view its details

After you pick an Compartment, the tool prints detailed information about the selected Compartment default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive images browser
   ocloud database compartment list
   ocloud database compartment list --json
`

// NewListCmd creates a new Cobra command for getting compartments in a specified tenancy or compartment.
// It supports pagination and optional JSON output.
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Compartments in the specified tenancy or compartment",
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

// RunListCommand handles the execution of the get command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running compartment list command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return compartment.ListCompartments(appCtx, useJSON)
}
