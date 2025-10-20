package cachecluster

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewCacheClusterCmd creates a new command for cluster cache-related operations
func NewCacheClusterCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "cache-cluster",
		Aliases:       []string{"cachecluster", "cc"},
		Short:         "Explore OCI Cache Clusters.",
		Long:          "Explore Oracle Cloud Infrastructure databases: list, get, and search",
		Example:       "  ocloud database cache-cluster list \n  ocloud database cache-cluster get \n  ocloud database cache-cluster search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))

	return cmd
}
