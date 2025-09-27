package loadbalancer

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

var findLong = ``

var findExamples = ``

func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find <pattern>",
		Short:         "Finds VCNs by a name pattern",
		Long:          findLong,
		Example:       findExamples,
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
	return nil
}
