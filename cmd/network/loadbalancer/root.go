package loadbalancer

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewLoadBalancerCmd creates a new command group for Load Balancer operations
func NewLoadBalancerCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "loadbalancer",
		Aliases:       []string{"lb", "lbr"},
		Short:         "Manage OCI Network Load Balancers",
		Long:          "Manage Oracle Cloud Infrastructure Network Load Balancers.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))
	return cmd
}
