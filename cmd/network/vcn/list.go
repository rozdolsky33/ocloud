package vcn

import (
	networkFlags "github.com/rozdolsky33/ocloud/cmd/network/flags"
	vcnFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
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

  # Include all related resources at once
  ocloud network vcn list --all

  # Output in JSON
  ocloud network vcn list --json

  # Using short aliases
  ocloud network vcn list -A -j
`

// NewListCmd creates a Cobra command named "list" that lists VCNs in a compartment.
// The returned command registers flags to include related network resources (gateway, subnet, NSG, route table, security list) and an all-info flag, silences usage/errors, and supports filtering and JSON output when executed.
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
	vcnFlags.AllInfoFlag.Add(cmd)

	return cmd
}

// runListCommand reads VCN-related flags from the provided command to determine output format and which related resources
// (gateways, subnets, NSGs, route tables, security lists) to include. If the `--all` flag is set, all related resources
// are included. It performs the VCN listing operation and returns any error encountered.
func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)
	gateways := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameGateway, false)
	subnets := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameSubnet, false)
	nsgs := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameNsg, false)
	routes := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameRoute, false)
	securityLists := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameSecurity, false)
	showAll := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameAll, false)

	if showAll {
		gateways, subnets, nsgs, routes, securityLists = true, true, true, true, true
	}

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn list", "json", useJSON, "all", showAll)

	return netvcn.ListVCNs(appCtx, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
