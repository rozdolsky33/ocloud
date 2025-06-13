package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/pkg/resources/compute"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// NewInstanceCmd creates a new command for instance-related operations
func NewInstanceCmd(appCtx *app.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "instance",
		Short:         "Find and list OCI instances",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInstanceCommand(cmd, appCtx)
		},
	}

	// Add all instance-related flags to the command
	flags.AddInstanceFlags(cmd)

	// Add subcommands
	cmd.AddCommand(
		newListCmd(appCtx),
		newFindCmd(appCtx),
	)

	return cmd
}

// doInstanceCommand handles the actual execution of instance commands based on config.
func doInstanceCommand(cmd *cobra.Command, appCtx *app.AppContext) error {
	list, _ := cmd.Flags().GetBool(flags.FlagNameList)
	find, _ := cmd.Flags().GetString(flags.FlagNameFind)
	imageDetails, _ := cmd.Flags().GetBool(flags.FlagNameImageDetails)

	switch {
	case list:
		fmt.Println("Listing instances in compartment:", appCtx.CompartmentName)
		return compute.ListInstances(appCtx)

	case find != "":
		fmt.Println("Finding instances with name pattern:", find)
		return compute.FindInstances(appCtx, find, imageDetails)

	default:
		return cmd.Help()
	}
}
