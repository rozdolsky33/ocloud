package networklb

import (
	nlbFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	nlbservice "github.com/rozdolsky33/ocloud/internal/services/network/networklb"
	"github.com/spf13/cobra"
)

var searchLong = `
Search for Network Load Balancers (L4) in the specified compartment that match the given pattern.

The search uses a combination of fuzzy, prefix, token, and substring matching across indexed fields.
You can search using any of the following fields (partial matches are supported):

Searchable fields:
- Name: Display name
- OCID: Network Load Balancer OCID
- Type: Public or private
- State: Lifecycle state
- VcnName: Name of the VCN
- IPAddresses: All assigned IP addresses
- Subnets: Subnet names/ids

Additional information:
- Use --all (-A) to include extra details in the output table
- Use --json (-j) to output the results in JSON format
- The search is case-insensitive. For highly specific inputs (like full OCIDs), exact and substring
  matches are attempted before broader fuzzy search.
`

var searchExamples = `
  # Search network load balancers whose name contains "prod"
  ocloud network nlb search prod

  # Search by IP address
  ocloud network nlb search 10.0.1

  # Include extra details in the table
  ocloud network nlb search prod --all

  # Use JSON output
  ocloud network nlb search prod --json

  # Short aliases
  ocloud net nlb s prod -A -j
`

func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search <pattern>",
		Aliases:       []string{"s"},
		Short:         "Fuzzy search for Network Load Balancers",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}

	nlbFlags.AllInfoFlag.Add(cmd)

	return cmd
}

func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	return nlbservice.SearchNetworkLoadBalancer(appCtx, namePattern, useJSON, showAll)
}
