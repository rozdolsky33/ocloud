package network

import (
	"github.com/rozdolsky33/ocloud/cmd/network/subnet"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewNetworkCmd creates a new cobra.Command for managing OCI network services such as vcn, subnets, load balancers and more
func NewNetworkCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "network",
		Aliases:       []string{"net"},
		Short:         "Manage OCI networking services",
		Long:          "Manage Oracle Cloud Infrastructure Networking services such as vcn, subnets and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(subnet.NewSubnetCmd(appCtx))

	return cmd
}
