package instance

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/pkg/resources"
)

// InstanceCmd is the root command for instance-related operations
var InstanceCmd = &cobra.Command{
	Use:           "instance",
	Short:         "Find and list OCI instances",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle the old flag syntax for backward compatibility
		list, _ := cmd.Flags().GetBool(config.FlagNameList)
		find, _ := cmd.Flags().GetString(config.FlagNameFind)
		imageDetails, _ := cmd.Flags().GetBool(config.FlagNameImageDetails)

		// Get the app context that was created in PersistentPreRunE
		ctx := cmd.Context()
		appCtx, ok := ctx.Value("appCtx").(*app.AppContext)
		if !ok {
			return fmt.Errorf("app context not found in command context")
		}

		if list {
			fmt.Println("Listing instances in compartment:", appCtx.CompartmentName)
			return resources.ListInstances(appCtx.Ctx, appCtx.Provider, appCtx.CompartmentID)
		} else if find != "" {
			fmt.Println("Finding instances with name pattern:", find)
			return resources.FindInstances(
				appCtx.Ctx,
				appCtx.Provider,
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

		// Initialize the logger
		if err := logger.SetLogger(); err != nil {
			return err
		}
		logger.InitLogger(logger.CmdLogger)

		// Create app context
		ctx := cmd.Context()
		appCtx, err := app.NewAppContext(ctx, cmd)
		if err != nil {
			return err
		}

		// Store app context in command context with a key
		ctx = context.WithValue(ctx, "appCtx", appCtx)
		cmd.SetContext(ctx)

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

func init() {
	InstanceCmd.Flags().BoolP(config.FlagNameList, config.FlagShortList, false, config.FlagDescList)
	InstanceCmd.Flags().StringP(config.FlagNameFind, config.FlagShortFind, "", config.FlagDescFind)
	InstanceCmd.Flags().BoolP(config.FlagNameImageDetails, config.FlagShortImageDetails, false, config.FlagDescImageDetails)
}
