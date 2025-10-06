package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewCompartmentCmd creates a new command for compartment-related operations
func NewCompartmentCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "compartment",
		Aliases:       []string{"compart", "comp", "cmp", "c"},
		Short:         "Manage OCI Compartments",
		Long:          "Manage Oracle Cloud Infrastructure Compartments: list, get and search",
		Example:       "  ocloud identity compartment get \n  ocloud identity compartment list \n  ocloud identity compartment search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))

	return cmd
}
