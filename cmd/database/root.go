package database

import (
	"github.com/rozdolsky33/ocloud/cmd/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/cmd/database/cachecluster"
	"github.com/rozdolsky33/ocloud/cmd/database/heatwave"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewDatabaseCmd creates a new cobra.Command to manage Oracle Cloud Infrastructure database services.
// It provides functionality for managing Autonomous Databases, HeatWave MySQL, and other database types.
func NewDatabaseCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "database",
		Aliases:       []string{"db"},
		Short:         "Explore OCI Database services",
		Long:          "Explore Oracle Cloud Infrastructure database services such as Autonomous Database, HeatWave and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(autonomousdb.NewAutonomousDatabaseCmd(appCtx))
	cmd.AddCommand(heatwave.NewHeatWaveDatabaseCmd(appCtx))
	cmd.AddCommand(cachecluster.NewCacheClusterCmd(appCtx))

	return cmd
}
