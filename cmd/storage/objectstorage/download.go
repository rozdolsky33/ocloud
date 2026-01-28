package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	osSvc "github.com/rozdolsky33/ocloud/internal/services/storage/objectstorage"
	"github.com/spf13/cobra"
)

var downloadLong = `
Download an object from an Object Storage bucket.

This command launches an interactive TUI that:
1. Shows available buckets to select the source
2. Shows objects in the bucket to select the file to download
3. Downloads the file to the current directory with progress reporting
`

var downloadExamples = `
  # Launch the interactive download flow
  ocloud storage object-storage download

  # Using short aliases
  ocloud stg os download
`

func NewDownloadCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "download",
		Short:         "Download an object from an Object Storage bucket",
		Long:          downloadLong,
		Example:       downloadExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return osSvc.DownloadFile(appCtx)
		},
	}

	return cmd
}
