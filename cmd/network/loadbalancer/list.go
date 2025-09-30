package loadbalancer

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	lbdomain "github.com/rozdolsky33/ocloud/internal/services/network/loadbalancer"
	"github.com/spf13/cobra"
)

var listLong = ``

var listExamples = ``

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
	return cmd
}

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := configflags.GetBoolFlag(cmd, configflags.FlagNameJSON, false)
	showAll := configflags.GetBoolFlag(cmd, configflags.FlagNameAll, false)
	return lbdomain.ListLoadBalancers(appCtx, useJSON, showAll)
}
