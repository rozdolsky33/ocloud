package networklb

import (
	nlbFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	nlbservice "github.com/rozdolsky33/ocloud/internal/services/network/networklb"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Network Load Balancers (L4) in the specified compartment using a TUI.

This command launches a terminal UI that loads available Network Load Balancers and lets you:
- Search/filter Network Load Balancers as you type
- Navigate the list
- Select a single Network Load Balancer to view its details

After you pick a Network Load Balancer, the tool prints detailed information about the selected NLB in the default table view or JSON format if specified with --json (-j).
You can also toggle inclusion of extra columns via --all (-A).
`

var listExamples = `
  # Launch the interactive Network Load Balancer browser
  ocloud network nlb list

  # Include extra columns in the table output
  ocloud network nlb list --all

  # Output in JSON
  ocloud network nlb list --json

  # Using short aliases
  ocloud net nlb list -A -j
`

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "Lists Network Load Balancers in a compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}
	nlbFlags.AllInfoFlag.Add(cmd)
	return cmd
}

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := configflags.GetBoolFlag(cmd, configflags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	return nlbservice.ListNetworkLoadBalancers(appCtx, useJSON, showAll)
}
