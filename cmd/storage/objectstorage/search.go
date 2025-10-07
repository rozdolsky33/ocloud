package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	objsvc "github.com/rozdolsky33/ocloud/internal/services/storage/objectstorage"
	"github.com/spf13/cobra"
)

var searchLong = `
Search for Object Storage Buckets in the specified compartment that match the given pattern.

The search uses a combination of fuzzy, prefix, token, and substring matching across indexed fields.
You can search using any of the following fields (partial matches are supported):

Searchable fields:
- Name: Bucket name
- OCID: Bucket OCID
- Namespace: Object storage namespace
- StorageTier: Standard/Archive, etc.
- Visibility: Public/Private
- Encryption: Encryption algorithm/provider
- Versioning: Versioning status
- TagsKV/TagsVal: Flattened tag keys and values
- ReplicationEnabled/IsReadOnly: Boolean flags

Additional information:
- Use --json (-j) to output the results in JSON format
- The search is case-insensitive. For highly specific inputs (like full OCIDs), exact and substring
  matches are attempted before broader fuzzy search.
`

var searchExamples = `
  # Buckets whose name contains "prod"
  ocloud storage objectstorage search prod

  # Search by namespace
  ocloud storage objectstorage search myns

  # Use JSON output
  ocloud storage objectstorage search prod --json

  # Short alias
  ocloud st os s prod -j
`

// NewSearchCmd creates a new command for searching buckets
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search <pattern>",
		Aliases:       []string{"s"},
		Short:         "Fuzzy search for Buckets",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunSearchCommand handles the execution of the search command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	pattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running object storage bucket search", "pattern", pattern, "json", useJSON)
	return objsvc.SearchBuckets(appCtx, pattern, useJSON)
}
