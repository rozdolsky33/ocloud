package instance

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/pkg/resources"
)

// newFindCmd creates a new command for finding instances by name
func newFindCmd(appCtx *app.AppContext) *cobra.Command {
	var showImageDetails bool

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

			return resources.FindInstances(appCtx, namePattern, showImageDetails)
		},
	}

	cmd.Flags().BoolVarP(&showImageDetails, "image-details", "i", false, "Show image details")

	return cmd
}
