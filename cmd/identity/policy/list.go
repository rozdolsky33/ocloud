package policy

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/policy"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Policies in the specified compartment using a TUI.

This command launches terminal UI that loads available Policies and lets you:
- Search/filter Policy as you type
- Navigate the list
- Select a single Policies to view its details

After you pick an Policies, the tool prints detailed information about the selected Policies default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive images browser
   ocloud identity policy list
   ocloud identity policy list --json
`

// NewListCmd returns "policy list".
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
	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running policy list command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return policy.ListPolicies(appCtx, useJSON)
}
