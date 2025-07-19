package auth

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// Short description for the auth command
var authShort = "Authenticate with OCI and refresh session tokens"

// Long description for the auth command
var authLong = `Provides commands for authenticating with Oracle Cloud Infrastructure (OCI).

This command group includes subcommands for authenticating with OCI and refreshing session tokens.
It allows you to interactively select your desired profile and region for authentication.`

// Examples for the auth command
var authExamples = `  ocloud config auth authenticate
  ocloud config auth authenticate -e
  ocloud config auth authenticate --filter us`

// NewAuthCmd creates a new cobra.Command for the auth command group.
func NewAuthCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "auth",
		Aliases:       []string{"authenticate", "a"},
		Short:         authShort,
		Long:          authLong,
		Example:       authExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is specified, run the authenticate command
			return NewAuthenticateCmd(appCtx).RunE(cmd, args)
		},
	}

	// Add subcommands
	cmd.AddCommand(NewAuthenticateCmd(appCtx))

	return cmd
}
