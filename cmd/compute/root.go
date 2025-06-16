package compute

import (
	"github.com/rozdolsky33/ocloud/cmd/compute/instance"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

// NewComputeCmd creates a new command for compute-related operations
func NewComputeCmd(appCtx *app.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "compute",
		Short:         "Manage OCI compute services",
		Long:          "Manage Oracle Cloud Infrastructure compute services such as instances, images, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add a custom help flag with a more descriptive message
	cmd.Flags().BoolP(flags.FlagNameHelp, flags.FlagShortHelp, false, flags.FlagDescHelp)
	_ = cmd.Flags().SetAnnotation(flags.FlagNameHelp, flags.CobraAnnotationKey, []string{"true"})

	// Add subcommands, passing in the AppContext
	cmd.AddCommand(instance.NewInstanceCmd(appCtx))

	return cmd
}
