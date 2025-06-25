package database

import (
	"github.com/rozdolsky33/ocloud/cmd/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

func NewDatabaseCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "database",
		Aliases:       []string{"db"},
		Short:         "Manage OCI Database services",
		Long:          "Manage Oracle Cloud Infrastructure database services such as Autonomous Database, HeatWave MySql and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands, passing in the ApplicationContext
	cmd.AddCommand(autonomousdb.NewAutonomousDatabaseCmd(appCtx))

	return cmd
}
