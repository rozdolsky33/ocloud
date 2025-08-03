package configuration

import (
	"github.com/rozdolsky33/ocloud/cmd/configuration/auth"
	"github.com/rozdolsky33/ocloud/cmd/configuration/info"
	"github.com/rozdolsky33/ocloud/cmd/configuration/setup"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the `configuration` command for managing ocloud CLI configurations, authentication with OCI,
// and viewing configuration information such as tenancy mappings.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "config",
		Aliases:       []string{"conf"},
		Short:         "Manage ocloud CLI configurations file and authentication",
		Long:          "Manage ocloud CLI configurations file and authentication with Oracle Cloud Infrastructure (OCI).\n\nThis command group provides functionality for:\n- Authenticating with OCI and refreshing session tokens\n- Viewing configuration information such as tenancy mappings",
		Example:       "  ocloud config session\n  ocloud config info\n ",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(info.NewInfoCmd())
	cmd.AddCommand(auth.NewSessionCmd())
	cmd.AddCommand(setup.SetupMappingFile())

	return cmd
}
