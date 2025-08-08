package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

func NewBastionCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "bastion",
		Aliases:       []string{"b"},
		Short:         "Manage OCI Bastion",
		Long:          "Manage Oracle Cloud Infrastructure Bastions - list created bastions or create new one.",
		Example:       "  ocloud identity bastion list \n  ocloud identity bastion create",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewCreateCmd(appCtx))

	return cmd
}
