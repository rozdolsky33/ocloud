package vcn

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
)

func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find <pattern>",
		Short:         "Finds VCNs by a name pattern",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFindCommand(cmd, args, appCtx)
		},
	}
	return cmd
}

func runFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	pattern := args[0]
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn find", "pattern", pattern, "json", useJSON)
	//netvcn.FindVCNs(appCtx, pattern, useJSON)
	return nil
}
