package instance

import (
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
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
	useJSON, _ := cmd.Flags().GetBool(flags.FlagNameJSON)

	switch {
	case list:
		// Get pagination parameters
		limit, err := cmd.Flags().GetInt(flags.FlagNameLimit)
		if err != nil {
			// Use default if flag not found
			limit = flags.FlagDefaultLimit
		}

		page, err := cmd.Flags().GetInt(flags.FlagNamePage)
		if err != nil {
			// Use default if flag not found
			page = flags.FlagDefaultPage
		}

		// Use VerboseInfo to ensure debug logs work with shorthand flags
		logger.VerboseInfo(logger.CmdLogger, 1, "Running instance list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON)
		return instance.ListInstances(appCtx, limit, page, useJSON)

	case find != "":
		// Use VerboseInfo to ensure debug logs work with shorthand flags
		logger.VerboseInfo(logger.CmdLogger, 1, "Running instance find command", "pattern", find, "in compartment", appCtx.CompartmentName, "json", useJSON)
		return instance.FindInstances(appCtx, find, imageDetails, useJSON)

	default:
		return cmd.Help()
	}
}
