package heatwave

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewHeatWaveDatabaseCmd creates a new command for database-related operations
func NewHeatWaveDatabaseCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "heatwave",
		Aliases:       []string{"hw"},
		Short:         "Explore OCI HeatWave Databases.",
		Long:          "Explore Oracle Cloud Infrastructure databases: list, get, and search",
		Example:       "  ocloud database heatwave list \n  ocloud database heatwave get \n  ocloud database heatwave search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))

	return cmd
}
