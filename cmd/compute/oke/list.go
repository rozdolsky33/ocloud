package oke

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	"github.com/spf13/cobra"
)

var listLong = `
List all Oracle Kubernetes Engine (OKE) clusters in the specified compartment.

This command displays information about all OKE clusters in the current compartment,
including their names, Kubernetes versions, endpoints, and associated node pools.
By default, it shows basic cluster information in a tabular format.

Additional Information:
- Use --json (-j) to output the results in JSON format
- Use --limit (-m) to control the number of results per page
- Use --page (-p) to navigate between pages of results
`

var listExamples = `
  # List all OKE clusters in the current compartment
  ocloud compute oke list

  # List all OKE clusters and output in JSON format
  ocloud compute oke list --json

  # List OKE clusters with pagination (10 per page, page 2)
  ocloud compute oke list --limit 10 --page 2
`

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Oracle Kubernetes Engine (OKE) clusters",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	// Add pagination flags
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, 1, "Running oke list command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return oke.ListClusters(appCtx, useJSON, limit, page)
}
