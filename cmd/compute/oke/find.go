package oke

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	"github.com/spf13/cobra"
)

var findLong = `
Find Oracle Kubernetes Engine (OKE) clusters in the specified compartment that match the given pattern.

This command searches for OKE clusters whose names or node pool names match the specified pattern.
By default, it shows detailed cluster information such as name, ID, Kubernetes version,
endpoint, and associated node pools for all matching clusters.

The search is performed using fuzzy matching, which means it will find clusters
even if the pattern is only partially matched. The search is case-insensitive.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command searches across all available clusters in the compartment
`

var findExamples = `
  # Find clusters with names containing "prod"
  ocloud compute oke find prod

  # Find clusters with names containing "dev" and output in JSON format
  ocloud compute oke find dev --json

  # Find clusters with names containing "test" (case-insensitive)
  ocloud compute oke find test
`

// NewFindCmd creates a new command for finding OKE clusters by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find OKE clusters by name pattern",
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
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running oke find command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return oke.FindClusters(appCtx, namePattern, useJSON)
}
