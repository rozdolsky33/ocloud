package auth

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/auth"
	"github.com/spf13/cobra"
)

// Short description for the authenticate command
var authenticateShort = "Authenticate with OCI and refresh session tokens"

// Long description for the authenticate command
var authenticateLong = `Runs the OCI CLI's session authenticate under the hood:

    oci session authenticate --profile-name <PROFILE> --region <REGION>

Interactively lets you pick your desired profile and region.`

// Examples for the authenticate command
var authenticateExamples = `  ocloud config auth authenticate
  ocloud config auth authenticate -e
  ocloud config auth authenticate --filter us`

// NewAuthenticateCmd creates a new cobra.Command for authenticating with OCI.
func NewAuthenticateCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "authenticate",
		Aliases:       []string{"auth", "a"},
		Short:         authenticateShort,
		Long:          authenticateLong,
		Example:       authenticateExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.LogWithLevel(appCtx.Logger, 1, "Initializing application")
			envOnly := flags.GetBoolFlag(cmd, flags.FlagNameEnvOnly, false)
			filter := flags.GetStringFlag(cmd, flags.FlagNameFilter, "")
			return auth.AuthenticateWithOCI(appCtx, envOnly, filter)
		},
	}

	// Add flags
	EnvOnlyFlag.Add(cmd)
	FilterFlag.Add(cmd)

	return cmd
}
