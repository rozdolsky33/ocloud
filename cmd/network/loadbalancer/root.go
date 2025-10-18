package loadbalancer

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewLoadBalancerCmd creates a new command group for Load Balancer operations
func NewLoadBalancerCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "load-balancer",
		Aliases:       []string{"loadbalancer", "lb", "lbr"},
		Short:         "Explore OCI Network Load Balancers",
		Long:          "Explore Oracle Cloud Infrastructure Network Load Balancers such as LBs, listeners, backend sets, and more",
		Example:       "  ocloud network load-balancer get \n  ocloud network load-balancer list \n  ocloud network load-balancer search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	return cmd
}
