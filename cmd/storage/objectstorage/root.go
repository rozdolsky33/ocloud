package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

var osLong = `
Manage Oracle Cloud Infrastructure Object Storage buckets and objects.

Commands:
  list      - Interactive TUI to browse buckets and objects, view details or download
  get       - Paginated listing of buckets with optional JSON output
  search    - Fuzzy search for buckets by name, tags, or other attributes
  upload    - Upload a file to a bucket (supports multipart for large files)
  download  - Download an object from a bucket
`

var osExamples = `
  # Browse buckets and objects interactively
  ocloud stg os list

  # List buckets with pagination
  ocloud stg os get --limit 10 --page 1

  # Search for buckets
  ocloud stg os search prod

  # Upload a file to a bucket (interactive)
  ocloud stg os upload

  # Download an object from a bucket (interactive)
  ocloud stg os download
`

func NewObjectStorageCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "object-storage",
		Aliases:       []string{"objectstorage", "os"},
		Short:         "Manage OCI Object Storage buckets and objects",
		Long:          osLong,
		Example:       osExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	cmd.AddCommand(NewUploadCmd(appCtx))
	cmd.AddCommand(NewDownloadCmd(appCtx))
	return cmd
}
