package objectstorage

import (
	osflags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	osSvc "github.com/rozdolsky33/ocloud/internal/services/storage/objectstorage"
	"github.com/spf13/cobra"
)

var getLong = `Get Object Storage buckets in the specified compartment with pagination support.

This command lists Object Storage buckets in the current compartment. By default, it shows a concise table
with key fields (name, namespace, created). Use --all (-A) to include extended bucket details (tier, access,
versioning, encryption, counts) and --json (-j) for machine-readable output.`

var getExamples = `  # Get buckets with default pagination (20 per page)
  ocloud storage object-storage get

  # Get buckets with custom pagination (10 per page, page 2)
  ocloud storage object-storage get --limit 10 --page 2

  # Get buckets and include extended details in the table
  ocloud storage object-storage get --all

  # JSON output with short aliases
  ocloud storage os get -A -j`

func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Paginated Object Storage Result",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}
	osflags.LimitFlag.Add(cmd)
	osflags.PageFlag.Add(cmd)
	osflags.AllInfoFlag.Add(cmd)
	return cmd
}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, osflags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, osflags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running object storage get command", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON, "all", showAll)
	return osSvc.GetBuckets(appCtx, limit, page, useJSON, showAll)
}
