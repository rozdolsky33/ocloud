package loadbalancer

import (
	lbFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	lbdomain "github.com/rozdolsky33/ocloud/internal/services/network/loadbalancer"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Load Balancers in the specified compartment using a TUI.

This command launches a terminal UI that loads available Load Balancers and lets you:
- Search/filter Load Balancers as you type
- Navigate the list
- Select a single Load Balancer to view its details

After you pick a Load Balancer, the tool prints detailed information about the selected Load Balancer in the default table view or JSON format if specified with --json (-j).
You can also toggle inclusion of extra columns via --all (-A).
`

var listExamples = `
  # Launch the interactive Load Balancer browser
  ocloud network loadbalancer list

  # Include extra columns in the table output
  ocloud network loadbalancer list --all

  # Output in JSON
  ocloud network loadbalancer list --json

  # Using short aliases
  ocloud net lb list -A -j
`

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "Lists Load Balancers in a compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}
	lbFlags.AllInfoFlag.Add(cmd)
	return cmd
}

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := configflags.GetBoolFlag(cmd, configflags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	return lbdomain.ListLoadBalancers(appCtx, useJSON, showAll)
}
