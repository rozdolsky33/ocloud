package cachecluster

import (
	cacheClusterFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/cacheclusterdb"
	"github.com/spf13/cobra"
)

var searchLong = `
Fuzzy Search for OCI Cache Clusters (Redis/Valkey) in the specified compartment.

Search across multiple cache cluster attributes including name, OCID, software version, cluster mode,
networking, and tags. The search uses fuzzy matching to find clusters even with typos or partial matches.

Searchable fields include:
  - Name, OCID, State, Lifecycle Details
  - Software Version (REDIS_7_0, VALKEY_7_2)
  - Cluster Mode (SHARDED, NONSHARDED)
  - Node Count, Node Memory Size, Shard Count
  - Configuration Set ID
  - VCN Name/ID, Subnet Name/ID
  - Network Security Group Names/IDs
  - Primary/Replicas/Discovery Endpoints (FQDN and IP addresses)
  - Tags (both keys and values)
`

var searchExamples = `
  # Search by cluster name
  ocloud database cache-cluster search prod-cache

  # Search by software version
  ocloud database cache-cluster search VALKEY_7_2

  # Search by cluster mode
  ocloud database cache-cluster search SHARDED

  # Search by VCN name
  ocloud database cache-cluster search prod-vcn

  # Search by endpoint FQDN
  ocloud database cache-cluster search cache.redis.us-ashburn-1

  # Search by IP address
  ocloud database cache-cluster search 10.0.20

  # Search by lifecycle state
  ocloud database cache-cluster search ACTIVE

  # Search by tag value
  ocloud database cache-cluster search production

  # Search with JSON output
  ocloud database cache-cluster search prod-cache --json

  # Search with detailed output
  ocloud database cache-cluster search prod-cache --all
`

// NewSearchCmd creates a new command for searching OCI Cache Clusters.
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for OCI Cache Clusters",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}
	cacheClusterFlags.AllInfoFlag.Add(cmd)
	return cmd
}

// runSearchCommand handles the execution of the search command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running OCI Cache Clusters search command", "searchPattern", namePattern, "json", useJSON, "showAll", showAll)
	return cacheclusterdb.SearchCacheClusters(appCtx, namePattern, useJSON, showAll)
}
