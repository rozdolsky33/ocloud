package compute

import (
	"github.com/rozdolsky33/ocloud/cmd/compute/instance"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// NewComputeCmd creates a new command for compute-related operations
func NewComputeCmd(appCtx *app.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "compute",
		Short:         "Manage OCI compute resources",
		Long:          "Manage Oracle Cloud Infrastructure compute resources such as instances, images, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add a custom help flag with a more descriptive message
	cmd.Flags().BoolP("help", "h", false, "help for compute (shorthand: -h)")
	_ = cmd.Flags().SetAnnotation("help", "cobra_annotation_flag_set_by_cobra", []string{"true"})

	// Add subcommands, passing in the AppContext
	cmd.AddCommand(instance.NewInstanceCmd(appCtx))

	return cmd
}
