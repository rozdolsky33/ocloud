package network

import (
	"github.com/rozdolsky33/ocloud/cmd/network/subnet"
	vcncmd "github.com/rozdolsky33/ocloud/cmd/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewNetworkCmd creates a Cobra command for managing OCI network services such as VCNs, subnets, and related resources.
// It configures the command (use "network", alias "net") and attaches the subnet and VCN subcommands before returning the constructed *cobra.Command.
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

	return cmd
}
