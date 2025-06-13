package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/pkg/resources/compute"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// NewInstanceCmd creates a new command for instance-related operations
func NewInstanceCmd(appCtx *app.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "instance",
		Short:         "Find and list OCI instances",
		Long:          "Find and list OCI instances using flags. Use --list or -l to list all instances, or --find or -f to find instances by name pattern.",
		Example:       "  ocloud compute instance --list\n  ocloud compute instance -l\n  ocloud compute instance --find myinstance\n  ocloud compute instance -f myinstance",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInstanceCommand(cmd, appCtx)
		},
	}

	// Add all instance-related flags to the command
	flags.AddInstanceFlags(cmd)

	// Add a custom help flag with a more descriptive message
	cmd.Flags().BoolP(flags.FlagNameHelp, flags.FlagShortHelp, false, flags.FlagDescHelp)
	_ = cmd.Flags().SetAnnotation(flags.FlagNameHelp, flags.CobraAnnotationKey, []string{"true"})

	return cmd
}

// doInstanceCommand handles the actual execution of instance commands based on config.
func doInstanceCommand(cmd *cobra.Command, appCtx *app.AppContext) error {
	list, _ := cmd.Flags().GetBool(flags.FlagNameList)
	find, _ := cmd.Flags().GetString(flags.FlagNameFind)
	imageDetails, _ := cmd.Flags().GetBool(flags.FlagNameImageDetails)

	switch {
	case list:
		// Use VerboseInfo to ensure debug logs work with shorthand flags
		logger.VerboseInfo(logger.CmdLogger, 1, "Running instance list command")
		fmt.Println("Listing instances in compartment:", appCtx.CompartmentName)
		return compute.ListInstances(appCtx)

	case find != "":
		// Use VerboseInfo to ensure debug logs work with shorthand flags
		logger.VerboseInfo(logger.CmdLogger, 1, "Running instance find command", "pattern", find)
		fmt.Println("Finding instances with name pattern:", find)
		return compute.FindInstances(appCtx, find, imageDetails)

	default:
		return cmd.Help()
	}
}
