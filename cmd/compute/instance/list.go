package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/compute/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new command for listing instances
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List all Instances",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	instaceFlags.ImageDetailsFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running instance list command in", "compartment", appCtx.CompartmentName, useJSON)
	return instance.ListInstances(appCtx, useJSON)
}
