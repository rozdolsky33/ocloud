package oke

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	"github.com/spf13/cobra"
)

// Dedicated documentation for the list command (separate from get)
var listLong = `
Interactively browse and search oke cluster in the specified compartment using a TUI.

This command launches terminal UI that loads available oke cluster and lets you:
- Search/filter oke cluster as you type
- Navigate the list
- Select a single oke cluster to view its details

After you pick an oke cluster, the tool prints detailed information about the selected oke cluster default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive oke cluster browser
  ocloud compute oke list
  ocloud compute oke list --json
`

// NewListCmd creates a new command for listing OKE clusters
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

	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running oke list command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return oke.ListClusters(appCtx, useJSON)
}
