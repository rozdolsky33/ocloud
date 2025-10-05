package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	configflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	osSvc "github.com/rozdolsky33/ocloud/internal/services/storage/objectstorage"
	"github.com/spf13/cobra"
)

var listLong = `
Interactively browse and search Object Storage Buckets in the specified compartment using a TUI.

This command launches a terminal UI that loads available Buckets and lets you:
- Search/filter Buckets as you type
- Navigate the list
- Select a single Bucket to view its details

After you pick a Bucket, the tool prints detailed information about the selected Bucket in the default table view or JSON format if specified with --json (-j).
`

var listExamples = `
  # Launch the interactive Bucket browser
  ocloud storage object-storage list

  # Output in JSON
  ocloud storage object-storage list --json

  # Using short aliases
  ocloud stg os list -j
`

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "Lists Object Storage Buckets in a compartment",
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

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := configflags.GetBoolFlag(cmd, configflags.FlagNameJSON, false)
	return osSvc.ListBuckets(appCtx, useJSON)
}
