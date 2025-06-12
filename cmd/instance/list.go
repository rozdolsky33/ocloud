package instance

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/pkg/resources"
)

// newListCmd creates a new command for listing instances
func newListCmd(appCtx *app.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List all instances in the compartment",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.CmdLogger.V(1).Info("Running instance list command")
			fmt.Println("Listing instances in compartment:", appCtx.CompartmentName)
			return resources.ListInstances(appCtx)
		},
	}

	return cmd
}
