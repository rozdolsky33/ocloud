package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/spf13/cobra"
)

// NewListCmd returns "bastion list".
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all bastions",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}
	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	ctx := cmd.Context()
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running list command")
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	err := bastion.ListBastions(ctx, appCtx, useJSON)
	if err != nil {
		return err
	}
	logger.CmdLogger.V(logger.Info).Info("Bastion list command completed.")
	return nil
}
