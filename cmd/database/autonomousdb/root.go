package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewAutonomousDatabaseCmd creates a new command for database-related operations
func NewAutonomousDatabaseCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "autonomous",
		Aliases:       []string{"adb"},
		Short:         "Manage OCI Databases.",
		Long:          "Manage Oracle Cloud Infrastructure databases: list, get, and search",
		Example:       "  ocloud database autonomous list \n  ocloud database autonomous get \n  ocloud database autonomous search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
