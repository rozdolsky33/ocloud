package vcn

import (
	networkFlags "github.com/rozdolsky33/ocloud/cmd/network/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	netvcn "github.com/rozdolsky33/ocloud/internal/services/network/vcn"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search VCNs in the specified compartment using a TUI.

This command launches a terminal UI that loads available Virtual Cloud Networks (VCNs) and lets you:
- Search/filter VCNs as you type
- Navigate the list
- Select a single VCN to view its details

After you pick a VCN, the tool prints detailed information about the selected VCN in the default table view or JSON format if specified with --json (-j).
You can also toggle inclusion of related networking resources via flags.
`

var listExamples = `
  # Launch the interactive VCN browser
  ocloud network vcn list

  # Launch and include related network resources
  ocloud network vcn list --gateway --subnet --nsg --route-table --security-list

  # Output in JSON
  ocloud network vcn list --json

  # Using short aliases
  ocloud network vcn list -G -S -N -R -L -j
`

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "Lists VCNs in a compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}

	networkFlags.Gateway.Add(cmd)
	networkFlags.Subnet.Add(cmd)
	networkFlags.Nsg.Add(cmd)
	networkFlags.RouteTable.Add(cmd)
	networkFlags.SecurityList.Add(cmd)

	return cmd
}

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)
	gateways := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameGateway, false)
	subnets := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameSubnet, false)
	nsgs := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameNsg, false)
	routes := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameRoute, false)
	securityLists := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameSecurity, false)

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn list", "json", useJSON)

	return netvcn.ListVCNs(appCtx, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
