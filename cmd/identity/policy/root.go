package policy

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewPolicyCmd creates a new command for policy-related operations
func NewPolicyCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "policy",
		Aliases:       []string{"pol"},
		Short:         "Manage OCI Policies",
		Long:          "Manage Oracle Cloud Infrastructure Policies - list all policies or find policy by pattern.",
		Example:       "  ocloud identity policy list \n  ocloud identity policy find mycompartment",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
