package networklb

import (
	nlbFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	nlbservice "github.com/rozdolsky33/ocloud/internal/services/network/networklb"
	"github.com/spf13/cobra"
)

var getLong = `Get all network load balancers (L4) in the specified compartment with pagination support.

This command displays information about network load balancers in the current compartment.
By default, it shows a concise table with key fields. Use flags to control pagination
and detail level.`

var getExamples = `  # Get all network load balancers with default pagination (20 per page)
  ocloud network nlb get

  # Get network load balancers with custom pagination (10 per page, page 2)
  ocloud network nlb get --limit 10 --page 2

  # Get network load balancers and include extra details in the table
  ocloud network nlb get --all

  # Output in JSON format
  ocloud net nlb get --json`

func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get Network Load Balancer Paginated Results",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}

	nlbFlags.LimitFlag.Add(cmd)
	nlbFlags.PageFlag.Add(cmd)
	nlbFlags.AllInfoFlag.Add(cmd)
	return cmd
}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := configflags.GetIntFlag(cmd, configflags.FlagNameLimit, nlbFlags.FlagDefaultLimit)
	page := configflags.GetIntFlag(cmd, configflags.FlagNamePage, nlbFlags.FlagDefaultPage)
	useJSON := configflags.GetBoolFlag(cmd, configflags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network load balancer get command", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON, "all", showAll)
	return nlbservice.GetNetworkLoadBalancers(appCtx, useJSON, limit, page, showAll)
}
