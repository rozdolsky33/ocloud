package vcn

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewVcnCmd creates and returns a Cobra command group for managing Oracle Cloud Infrastructure Virtual Cloud Networks (VCNs).
// The command is configured with usage "vcn", descriptive short and long help text, silenced usage and errors, and it registers
// the get, list, and find subcommands using the provided application context.
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
	cmd.AddCommand(NewFindCmd(appCtx))
	return cmd
}
