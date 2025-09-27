package loadbalancer

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

var getLong = ``

var getExamples = ``

func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get Load Balancers",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(cmd, appCtx)
		},
	}
	return cmd
}

func runGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {

	return nil
}
