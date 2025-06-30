package identity

import (
	"github.com/rozdolsky33/ocloud/cmd/identity/compartment"
	"github.com/rozdolsky33/ocloud/cmd/identity/policy"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewIdentityCmd creates a new cobra.Command for managing OCI identity services such as compartments, polices and bastions.
func NewIdentityCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "identity",
		Aliases:       []string{"ident", "idt"},
		Short:         "Manage OCI identity services",
		Long:          "Manage Oracle Cloud Infrastructure Identity services such as compartments, policies, bastion sessions and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands, passing in the ApplicationContext
	cmd.AddCommand(compartment.NewCompartmentCmd(appCtx))
	cmd.AddCommand(policy.NewPolicyCmd(appCtx))

	return cmd
}
