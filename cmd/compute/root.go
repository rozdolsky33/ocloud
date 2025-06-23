package compute

import (
	"github.com/rozdolsky33/ocloud/cmd/compute/image"
	"github.com/rozdolsky33/ocloud/cmd/compute/instance"
	"github.com/rozdolsky33/ocloud/cmd/compute/oke"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// NewComputeCmd creates a new command for compute-related operations
func NewComputeCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "compute",
		Aliases:       []string{"comp"},
		Short:         "Manage OCI compute services",
		Long:          "Manage Oracle Cloud Infrastructure Compute services such as instances, image, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands, passing in the ApplicationContext
	cmd.AddCommand(instance.NewInstanceCmd(appCtx))
	cmd.AddCommand(image.NewImageCmd(appCtx))
	cmd.AddCommand(oke.NewOKECmd(appCtx))

	return cmd
}
