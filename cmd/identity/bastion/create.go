package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
)

func NewCreateCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create",
		Aliases:       []string{"c"},
		Short:         "",
		Long:          "",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCreateCommand(cmd, appCtx)
		},
	}
	// Add flags specific to the list command

	return cmd

}

func RunCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Get pagination parameters

	logger.LogWithLevel(logger.CmdLogger, 1, "Running bastion list command in", "compartment", appCtx.CompartmentName)
	return nil
}
