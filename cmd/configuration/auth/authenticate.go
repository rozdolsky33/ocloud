package auth

import (
	configurationFlags "github.com/rozdolsky33/ocloud/cmd/configuration/flags"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/auth"
	"github.com/spf13/cobra"
)

// Short description for the authenticate command
var authenticateShort = "Authenticate with OCI and refresh session tokens"

// Long description for the authenticate command
var authenticateLong = `Interactively guides you through the authentication process with OCI.
Allows you to select your desired profile and region.
You can use --filter to filter regions by prefix and --realm to filter by realm.
If a tenancy-mapping file is present, the --realm flag will also filter tenancy mappings by the specified realm.`

// Examples for the authenticate command
var authenticateExamples = `  ocloud config session authenticate 
  ocloud config session authenticate --filter us
  ocloud config session authenticate --realm OC1
  ocloud config session auth -f us -r OC2`

// NewAuthenticateCmd creates a new cobra.Command for authenticating with OCI.
func NewAuthenticateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "authenticate",
		Aliases:       []string{"auth", "a"},
		Short:         authenticateShort,
		Long:          authenticateLong,
		Example:       authenticateExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunAuthenticateCommand(cmd)
		},
	}
	// Add filter flag
	configurationFlags.FilterFlag.Add(cmd)

	// Add realm filter flag
	configurationFlags.RealmFlag.Add(cmd)

	return cmd
}

func RunAuthenticateCommand(cmd *cobra.Command) error {
	filter := flags.GetStringFlag(cmd, flags.FlagNameFilter, "")
	realm := flags.GetStringFlag(cmd, flags.FlagNameRealm, "")
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running authenticate command", "filter", filter, "realm", realm)
	err := auth.AuthenticateWithOCI(filter, realm)
	if err != nil {
		return err
	}
	logger.CmdLogger.V(logger.Info).Info("Authentication command completed.")
	return nil
}
