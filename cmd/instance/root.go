package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/pkg/resources"
)

// InstanceCmd is the root command for instance-related operations
var InstanceCmd = &cobra.Command{
	Use:     "instance",
	Short:   "Find and list OCI instances",
	PreRunE: preConfigE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle the old flag syntax for backward compatibility
		list, _ := cmd.Flags().GetBool(config.FlagNameList)
		find, _ := cmd.Flags().GetString(config.FlagNameFind)
		imageDetails, _ := cmd.Flags().GetBool(config.FlagNameImageDetails)

		if list {
			// Create app context
			ctx := cmd.Context()
			appCtx, err := app.NewAppContext(ctx, cmd)
			if err != nil {
				return err
			}

			fmt.Println("Listing instances in compartment:", appCtx.CompartmentName)
			return resources.ListInstances(appCtx.Ctx, appCtx.CompartmentID)
		} else if find != "" {
			// Create app context
			ctx := cmd.Context()
			appCtx, err := app.NewAppContext(ctx, cmd)
			if err != nil {
				return err
			}

			fmt.Println("Finding instances with name pattern:", find)
			return resources.FindInstances(
				appCtx.Ctx,
				appCtx.CompartmentID,
				find,
				imageDetails,
			)
		}

		// If no flags are provided, show help
		return cmd.Help()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip if it's the help command
		if cmd.Name() == "help" {
			return nil
		}

		// Create app context
		ctx := cmd.Context()
		appCtx, err := app.NewAppContext(ctx, cmd)
		if err != nil {
			return err
		}

		// Store app context in command
		cmd.SetContext(appCtx.Ctx)

		// Add subcommands
		if cmd.Name() == "instance" && len(cmd.Commands()) == 0 {
			cmd.AddCommand(
				newListCmd(appCtx),
				newFindCmd(appCtx),
			)
		}

		return nil
	},
}

// preConfigE initializes the logger
func preConfigE(cmd *cobra.Command, args []string) error {
	if err := logger.SetLogger(); err != nil {
		return err
	}
	logger.InitLogger(logger.CmdLogger)

	return nil
}

func init() {
	InstanceCmd.Flags().BoolP(config.FlagNameList, config.FlagShortList, false, config.FlagDescList)
	InstanceCmd.Flags().StringP(config.FlagNameFind, config.FlagShortFind, "", config.FlagDescFind)
	InstanceCmd.Flags().BoolP(config.FlagNameImageDetails, config.FlagShortImageDetails, false, config.FlagDescImageDetails)
}
