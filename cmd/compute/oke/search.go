package oke

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	"github.com/spf13/cobra"
)

var findLong = `
FuzzySearch Oracle Kubernetes Engine (OKE) clusters in the specified compartment that match the given pattern.

This command searches across both cluster-level and node-pool attributes. It matches on:
- Cluster: name, OCID, Kubernetes version, state, VCN OCID, private/public endpoints, and tags
- Node pools: display names and node shapes (aggregated per cluster)

By default, it shows detailed cluster information such as name, ID, Kubernetes version,
endpoints, and associated node pools for all matching clusters.

The search is performed using fuzzy matching, which means it will find clusters
even if the pattern is only partially matched. The search is case-insensitive.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command searches across all available clusters in the compartment
`

var findExamples = `
  # Fuzzy search clusters with names containing "prod"
  ocloud compute oke search prod

  # Fuzzy search clusters with names containing "dev" and output in JSON format
  ocloud compute oke search dev --json

  # Fuzzy search clusters by node pool shape
  ocloud compute oke search "VM.Standard3"

  # Fuzzy search clusters by OCID fragment or exact OCID
  ocloud compute oke search ocid1.clusters
  ocloud compute oke search ocid1.cluster.oc1..exampleexactocid

  # Fuzzy search clusters with tags (key or value)
  ocloud compute oke search team:platform
  ocloud compute oke search platform
`

// NewFindCmd creates a new command for finding OKE clusters by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for OKE clusters",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running oke search command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return oke.SearchOKEClusters(appCtx, namePattern, useJSON)
}
