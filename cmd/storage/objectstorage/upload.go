package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	osSvc "github.com/rozdolsky33/ocloud/internal/services/storage/objectstorage"
	"github.com/spf13/cobra"
)

var uploadLong = `
Upload a file to an Object Storage bucket using multipart upload for large files.

This command launches an interactive TUI that:
1. Shows available buckets to select the destination
2. Shows a file browser to select the file to upload
3. Uploads the file with progress reporting

Files larger than 10 MiB are automatically uploaded using multipart upload
for better performance and reliability.
`

var uploadExamples = `
  # Launch the interactive upload flow
  ocloud storage object-storage upload

  # Using short aliases
  ocloud stg os upload
`

func NewUploadCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "upload",
		Short:         "Upload a file to an Object Storage bucket",
		Long:          uploadLong,
		Example:       uploadExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return osSvc.UploadFile(appCtx)
		},
	}

	return cmd
}
