package oke

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	"github.com/spf13/cobra"
)

var getLong = `
Get all Oracle Kubernetes Engine (OKE) clusters in the specified compartment.

This command displays information about all OKE clusters in the current compartment,
including their names, Kubernetes versions, endpoints, and associated node pools.
By default, it shows basic cluster information in a tabular format.

Additional Information:
- Use --json (-j) to output the results in JSON format
- Use --limit (-m) to control the number of results per page
- Use --page (-p) to navigate between pages of results
`

var getExamples = `
  # Get all OKE clusters in the current compartment
  ocloud compute oke get

  # Get all OKE clusters and output in JSON format
  ocloud compute oke get --json

  # Get OKE clusters with pagination (10 per page, page 2)
  ocloud compute oke get --limit 10 --page 2
`

// NewGetCmd creates a new cobra.Command for listing all OKE clusters in a specified compartment.
// The command supports pagination through the --limit and --page flags for controlling list size and navigation.
// The command supports JSON output through the --JSON flag.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all Oracle Kubernetes Engine (OKE) clusters",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, appCtx)
		},
	}

	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd
}

// RunGetCommand handles the execution of the list command
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running oke get command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return oke.GetClusters(appCtx, useJSON, limit, page)
}
