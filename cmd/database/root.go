package database

import (
	"github.com/rozdolsky33/ocloud/cmd/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewDatabaseCmd creates a new cobra.Command to manage Oracle Cloud Infrastructure database services.
// It provides functionality for managing Autonomous Databases, HeatWave MySQL, and other database types.
func NewDatabaseCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "database",
		Aliases:       []string{"db"},
		Short:         "Manage OCI Database services",
		Long:          "Manage Oracle Cloud Infrastructure database services such as Autonomous Database, HeatWave MySql and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(autonomousdb.NewAutonomousDatabaseCmd(appCtx))

	return cmd
}
