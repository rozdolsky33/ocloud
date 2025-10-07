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

var searchLong = `
Search for Virtual Cloud Networks (VCNs) in the specified compartment that match the given pattern.

The search uses a combination of fuzzy, prefix, token, and substring matching across indexed fields.
You can search using any of the following fields (partial matches are supported):

Searchable fields:
- Name: Display name
- OCID: VCN OCID
- State: Lifecycle state
- CIDRs: IPv4/IPv6 CIDR blocks
- DnsLabel: DNS label
- DomainName: VCN domain name
- TagsKV/TagsVal: Flattened tag keys and values
- Gateways/Subnets/NSGs/RouteTables/SecLists: Related resource names

Additional information:
- Use --all (-A) to include related resources in the output (gateways, subnets, NSGs, route tables, security lists)
- Use --json (-j) to output the results in JSON format
- The search is case-insensitive. For highly specific inputs (like full OCIDs), exact and substring
  matches are attempted before broader fuzzy search.
`

var searchExamples = `
  # Search VCNs whose name contains "prod"
  ocloud network vcn search prod

  # Search by DNS label or domain name
  ocloud network vcn search corp

  # Include related resources in the output table
  ocloud network vcn search prod --all

  # Use JSON output
  ocloud network vcn search prod --json

  # Short aliases
  ocloud net vcn s prod -A -j
`

func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search <pattern>",
		Aliases:       []string{"s"},
		Short:         "Fuzzy search for VCNs",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
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

func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	pattern := args[0]
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
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn search", "pattern", pattern, "json", useJSON, "all", showAll)
	return netvcn.SearchVCNs(appCtx, pattern, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
