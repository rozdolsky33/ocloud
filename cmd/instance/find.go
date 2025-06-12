package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/pkg/flags"
	"github.com/rozdolsky33/ocloud/pkg/resources/compute"

	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// newFindCmd creates a new command for finding instances by name
func newFindCmd(appCtx *app.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [name]",
		Short:         "Find instances by name pattern",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			namePattern := args[0]
			logger.CmdLogger.V(1).Info("Running instance find command", "pattern", namePattern)
			fmt.Println("Finding instances with name pattern:", namePattern)

			showImageDetails, _ := cmd.Flags().GetBool(flags.FlagNameImageDetails)
			return compute.FindInstances(appCtx, namePattern, showImageDetails)
		},
	}

	flags.ImageDetailsFlag.AddBoolFlag(cmd)

	return cmd
}
