package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewBastionCmd creates a new command for compartment-related operations
func NewBastionCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "bastion",
		Aliases:       []string{"bast"},
		Short:         "Manage OCI Bastion Sessions",
		Long:          "Manage Oracle Cloud Infrastructure Bastion - list bastion session find  or create bastion.",
		Example:       "  ocloud identity bastion list \n  ocloud identity bastion find mycompartment",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands

	return cmd
}
