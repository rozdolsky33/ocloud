package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewCompartmentCmd creates a new command for compartment-related operations
func NewCompartmentCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "compartment",
		Aliases:       []string{"compart"},
		Short:         "Manage OCI Compartments",
		Long:          "Manage Oracle Cloud Infrastructure Compartments - list all compartments or find compartment by pattern.",
		Example:       "  ocloud identity compartment list \n  ocloud identity compartment find mycompartment",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
