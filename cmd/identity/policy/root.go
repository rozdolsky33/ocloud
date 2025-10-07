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
		Long:          "Manage Oracle Cloud Infrastructure Policies: list, get, and search",
		Example:       "  ocloud identity policy get \n  ocloud identity policy list \n  ocloud identity policy search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))

	return cmd
}
