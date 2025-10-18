package vcn

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewVcnCmd creates a new command group for VCN-related operations
func NewVcnCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vcn",
		Short:         "Explore OCI Virtual Cloud Networks (VCNs)",
		Long:          "  ocloud network vcn list \n  ocloud network vcn get \n  ocloud network vcn search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	return cmd
}
