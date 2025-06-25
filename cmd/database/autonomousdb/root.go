package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewAutonomousDatabaseCmd creates a new command for database-related operations
func NewAutonomousDatabaseCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "autonomous",
		Aliases:       []string{"auto"},
		Short:         "Manage OCI Compartments",
		Long:          "Manage Oracle Cloud Infrastructure Databases - list all databases or find database by pattern.",
		Example:       "  ocloud database autonomous list \n  ocloud database autonomous find mydatabase",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
