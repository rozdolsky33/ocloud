package configuration

import (
	"github.com/rozdolsky33/ocloud/cmd/configuration/auth"
	"github.com/rozdolsky33/ocloud/cmd/configuration/info"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the `configuration` command for managing ocloud CLI configurations and related operations.
func NewConfigCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "config",
		Aliases:       []string{"conf"},
		Short:         "Manage ocloud CLI configurations",
		Long:          "Manage ocloud CLI configurations with OCI such as authentication, view configuration information, and more.",
		Example:       "  ocloud config info map-file\n  ocloud config info map-file --json\n  ocloud config info map-file --realm OC1",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(info.NewInfoCmd(appCtx))
	cmd.AddCommand(auth.NewAuthCmd(appCtx))

	return cmd
}
