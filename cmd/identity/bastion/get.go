package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/spf13/cobra"
)

// NewGetCmd creates a new cobra.Command for listing bastions.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get all bastions",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}
	return cmd
}

// RunGetCommand handles the execution of the get command
func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	ctx := cmd.Context()
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running Get command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	return bastion.GetBastions(ctx, appCtx, useJSON)
}
