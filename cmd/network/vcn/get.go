package vcn

import (
	networkFlags "github.com/rozdolsky33/ocloud/cmd/network/flags"
	vcnFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	netvcn "github.com/rozdolsky33/ocloud/internal/services/network/vcn"
	"github.com/spf13/cobra"
)

// Long description for the get command
var getLong = `
Fetch VCNs in the specified compartment with pagination support.

This command retrieves Virtual Cloud Networks (VCNs) in the current compartment.
By default, it shows basic information such as name, OCID, state, compartment, and CIDR blocks.

The output is paginated. Control the number of VCNs per page with --limit (-m) and
navigate pages using --page (-p).

Additional Information:
- Use --json (-j) to output the results in JSON format
- Use flags to include related resources: gateways, subnets, NSGs, route tables, security lists
`

// Examples for the get command
var getExamples = `
  # Get VCNs with default pagination
  ocloud network vcn get

  # Get VCNs with custom pagination (10 per page, page 2)
  ocloud network vcn get --limit 10 --page 2

  # Include related resources
  ocloud network vcn get --gateway --subnet --nsg --route-table --security-list

  # Include all related resources at once
  ocloud network vcn get --all

  # JSON output with short aliases
  ocloud network vcn get -m 5 -p 3 -A -j
`

// NewGetCmd creates the "vcn get" Cobra command for fetching Virtual Cloud Networks (VCNs).
// The command supports pagination and optional inclusion of related resources (gateways, subnets,
// NSGs, route tables, security lists) and registers flags for those options as well as `--all`,
// `--limit`, and `--page`. The command's execution is delegated to runGetCommand with the
// provided application context.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get VCNs",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}

	networkFlags.Gateway.Add(cmd)
	networkFlags.Subnet.Add(cmd)
	networkFlags.Nsg.Add(cmd)
	networkFlags.RouteTable.Add(cmd)
	networkFlags.SecurityList.Add(cmd)
	vcnFlags.AllInfoFlag.Add(cmd)
	vcnFlags.LimitFlag.Add(cmd)
	vcnFlags.PageFlag.Add(cmd)

	return cmd
}

// runGetCommand runs the "vcn get" command, reading flags for pagination, JSON output, and related-resource inclusion (or enabling all with --all), then invokes netvcn.GetVCNs.
// It returns any error produced while retrieving or printing VCNs.
func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, vcnFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, vcnFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	gateways := flags.GetBoolFlag(cmd, flags.FlagNameGateway, false)
	subnets := flags.GetBoolFlag(cmd, flags.FlagNameSubnet, false)
	nsgs := flags.GetBoolFlag(cmd, flags.FlagNameNsg, false)
	routes := flags.GetBoolFlag(cmd, flags.FlagNameRoute, false)
	securityLists := flags.GetBoolFlag(cmd, flags.FlagNameSecurity, false)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	if showAll {
		gateways, subnets, nsgs, routes, securityLists = true, true, true, true, true
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn get", "json", useJSON, "all", showAll)
	return netvcn.GetVCNs(appCtx, limit, page, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
