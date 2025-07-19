package auth

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/auth"
	"github.com/spf13/cobra"
)

// NewAuthCmd creates a new cobra.Command for authenticating with OCI.
func NewAuthCmd(appCtx *app.ApplicationContext) *cobra.Command {
	var envOnly bool
	var filter string

	cmd := &cobra.Command{
		Use:     "auth",
		Aliases: []string{"a"},
		Short:   "Authenticate with OCI and refresh session tokens",
		Long: `Runs the OCI CLI's session authenticate under the hood:

    oci session authenticate --profile-name <PROFILE> --region <REGION>

Interactively lets you pick your desired profile and region.`,
		Example:       "  ocloud config auth\n  ocloud config auth -e\n  ocloud config auth --filter us",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.LogWithLevel(appCtx.Logger, 1, "Initializing application")
			return auth.AuthenticateWithOCI(appCtx, envOnly, filter)
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&envOnly, "env-only", "e", false, "Only output environment variables, don't run interactive authentication")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter regions by prefix (e.g., us, eu, ap)")

	return cmd
}
