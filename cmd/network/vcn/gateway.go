package vcn

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
)

// NewGetGatewayCmd creates a new command for finding subnets by name pattern
func NewGetGatewayCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get [vcn name]",
		Aliases:       []string{"f"},
		Short:         "Get VCN Gateway with vcn name",
		Long:          "",
		Example:       "",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunFindCommand handles the execution of the find command
func runCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	fmt.Println("Running find command")
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running gateways find command", "pattern", args[0])
	vcnName := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running gateways get command", "pattern", vcnName, "json", useJSON)
	return nil
}
