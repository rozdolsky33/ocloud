package network

import (
	lbcmd "github.com/rozdolsky33/ocloud/cmd/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/cmd/network/subnet"
	vcncmd "github.com/rozdolsky33/ocloud/cmd/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewNetworkCmd creates a new cobra.Command for managing OCI network services such as vcn, subnets, load balancers and more
func NewNetworkCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "network",
		Aliases:       []string{"net"},
		Short:         "Manage OCI network services",
		Long:          "Manage Oracle Cloud Infrastructure Networking services such as vcn, subnets and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(subnet.NewSubnetCmd(appCtx))
	cmd.AddCommand(vcncmd.NewVcnCmd(appCtx))
	cmd.AddCommand(lbcmd.NewLoadBalancerCmd(appCtx))

	return cmd
}
