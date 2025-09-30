package loadbalancer

import (
	lbFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	lbservice "github.com/rozdolsky33/ocloud/internal/services/network/loadbalancer"
	"github.com/spf13/cobra"
)

var getLong = `Get all load balancers in the specified compartment with pagination support.

This command displays information about load balancers in the current compartment.
By default, it shows a concise table with key fields. Use flags to control pagination
and detail level.`

var getExamples = `  # Get all load balancers with default pagination (20 per page)
  ocloud network loadbalancer get

  # Get load balancers with custom pagination (10 per page, page 2)
  ocloud network loadbalancer get --limit 10 --page 2

  # Get load balancers and include extra details in the table
  ocloud network loadbalancer get --all

  # Output in JSON format
  ocloud network loadbalancer get --json`

func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         " Get Load Balancer Paginated Results",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}

	lbFlags.LimitFlag.Add(cmd)
	lbFlags.PageFlag.Add(cmd)
	lbFlags.AllInfoFlag.Add(cmd)
	return cmd
}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := configflags.GetIntFlag(cmd, configflags.FlagNameLimit, lbFlags.FlagDefaultLimit)
	page := configflags.GetIntFlag(cmd, configflags.FlagNamePage, lbFlags.FlagDefaultPage)
	useJSON := configflags.GetBoolFlag(cmd, configflags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running load balancer get command", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON, "all", showAll)
	return lbservice.GetLoadBalancers(appCtx, useJSON, limit, page, showAll)
}
