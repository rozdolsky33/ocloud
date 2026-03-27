package networklb

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewNetworkLoadBalancerCmd creates a new command group for Network Load Balancer operations
func NewNetworkLoadBalancerCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "network-load-balancer",
		Aliases:       []string{"networkloadbalancer", "nlb"},
		Short:         "Explore OCI Network Load Balancers (L4)",
		Long:          "Explore Oracle Cloud Infrastructure Network Load Balancers (Layer 4) such as NLBs, listeners, backend sets, and more",
		Example:       "  ocloud network network-load-balancer get \n  ocloud network nlb list \n  ocloud network nlb search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	return cmd
}
