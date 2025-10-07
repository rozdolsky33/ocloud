package vcn

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewVcnCmd creates a new command group for VCN-related operations
func NewVcnCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vcn",
		Short:         "Manage OCI Virtual Cloud Networks (VCNs)",
		Long:          "Manage Oracle Cloud Infrastructure Virtual Cloud Networks (VCNs).",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	return cmd
}
