package objectstorage

import (
	osflags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	osSvc "github.com/rozdolsky33/ocloud/internal/services/storage/objectstorage"
	"github.com/spf13/cobra"
)

var getLong = ``

var getExamples = ``

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
	return cmd
}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, osflags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, osflags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running image list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON)
	return osSvc.GetBuckets(appCtx, limit, page, useJSON)
}
