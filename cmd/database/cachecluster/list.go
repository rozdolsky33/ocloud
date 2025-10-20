package cachecluster

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/cacheclusterdb"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search OCI Cache Clusters in the specified compartment using a TUI.

This command launches terminal UI that loads available OCI Cache Cluster and lets you:
- Search/filter OCI Cache Cluster  as you type
- Navigate the list
- Select a single OCI Cache Cluster to view its details

After you pick an OCI Cache Cluster, the tool prints detailed information about the selected OCI Cache Cluster default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive images browser
   ocloud database cache-cluster list
   ocloud database cache-cluster list --json
`

// NewListCmd creates a new command for listing OCI Cache Cluster
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all OCI Cache Clusters",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}
	return cmd

}

// runListCommand handles the execution of the list command
func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running OCI Cache Cluster list command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	return cacheclusterdb.ListCacheClusters(appCtx, useJSON)
}
