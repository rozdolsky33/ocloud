package cachecluster

import (
	cacheClusterFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/cacheclusterdb"
	"github.com/spf13/cobra"
)

// Long description for the list command
var getLong = `
Fetch OCI Cache Clusters in the specified compartment with pagination support.

This command displays information about available OCI Cache Clusters in the current compartment.
By default, it shows basic information such as name, ID, state, and workload type.

The output is paginated, with a default limit of 20 per page. You can navigate
through pages using the --page flag and control the number of per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available OCI Cache Clusters in the compartment
`

// Examples for the list command
var getExamples = `
  # Get all OCI Cache Clusters with default pagination (20 per page)
  ocloud database cache-cluster get

  # Get OCI Cache Clusters with custom pagination (10 per page, page 2)
  ocloud database cache-cluster get --limit 10 --page 2

  # Get OCI Cache Clusters and output in JSON format
  ocloud database cache-cluster get --json

  # Get OCI Cache Clusters with custom pagination and JSON output
  ocloud database cache-cluster get --limit 5 --page 3 --json
`

// NewGetCmd creates a "list" subcommand for listing all in the specified compartment with pagination support.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all OCI Cache Clusters",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}

	cacheClusterFlags.LimitFlag.Add(cmd)
	cacheClusterFlags.PageFlag.Add(cmd)
	cacheClusterFlags.AllInfoFlag.Add(cmd)

	return cmd

}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running CacheClusters database Get command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, cacheClusterFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, cacheClusterFlags.FlagDefaultPage)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	return cacheclusterdb.GetCacheClusters(appCtx, useJSON, limit, page, showAll)
}
