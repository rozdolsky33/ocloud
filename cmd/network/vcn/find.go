package vcn

import (
	networkFlags "github.com/rozdolsky33/ocloud/cmd/network/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	netvcn "github.com/rozdolsky33/ocloud/internal/services/network/vcn"
	"github.com/spf13/cobra"
)

var findLong = `
Find VCNs by display name using a pattern.

This command searches Virtual Cloud Networks (VCNs) in the current compartment using a case-insensitive
substring match against the VCN display name. By default, it prints a concise table of matches.

Use --json (-j) to output raw JSON. You can include related networking resources with flags.
`

var findExamples = `
  # Find VCNs whose name contains "prod"
  ocloud network vcn find prod

  # Include related resources (gateways, subnets, NSGs, route tables, security lists)
  ocloud network vcn find prod --gateway --subnet --nsg --route-table --security-list

  # JSON output with short aliases
  ocloud network vcn find prod -G -S -N -R -L -j
`

func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find <pattern>",
		Short:         "Finds VCNs by a name pattern",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFindCommand(cmd, args, appCtx)
		},
	}
	networkFlags.Gateway.Add(cmd)
	networkFlags.Subnet.Add(cmd)
	networkFlags.Nsg.Add(cmd)
	networkFlags.RouteTable.Add(cmd)
	networkFlags.SecurityList.Add(cmd)
	return cmd
}

func runFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	pattern := args[0]
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)
	gateways := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameGateway, false)
	subnets := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameSubnet, false)
	nsgs := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameNsg, false)
	routes := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameRoute, false)
	securityLists := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameSecurity, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn find", "pattern", pattern, "json", useJSON)
	return netvcn.FindVCNs(appCtx, pattern, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
