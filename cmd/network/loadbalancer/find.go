package loadbalancer

import (
	lbFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	lbservice "github.com/rozdolsky33/ocloud/internal/services/network/loadbalancer"
	"github.com/spf13/cobra"
)

var findLong = `
Find Load Balancers by display name using a pattern.

This command searches Load Balancers in the current compartment using a case-insensitive
substring and fuzzy match against the Load Balancer display name. By default, it prints a concise
table of matches. Use --all to include extra columns, and --json to output machine-readable JSON.
`

var findExamples = `
  # Find load balancers whose name contains "prod"
  ocloud network loadbalancer find prod

  # Use JSON output
  ocloud network loadbalancer find prod --json

  # Include extra details in the table
  ocloud network loadbalancer find prod --all

  # Short aliases
  ocloud net lb find prod -A -j
`

func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find <pattern>",
		Aliases:       []string{"f"},
		Short:         "Finds Load Balancer with existing attribute",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFindCommand(cmd, args, appCtx)
		},
	}

	lbFlags.AllInfoFlag.Add(cmd)

	return cmd
}

func runFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	return lbservice.FindLoadBalancer(appCtx, namePattern, useJSON, showAll)
}
