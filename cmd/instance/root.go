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
	PreRunE:       setupInstanceContext,
	RunE:          executeInstanceCommand,
}

func init() {
	InstanceCmd.Flags().BoolP(config.FlagNameList, config.FlagShortList, false, config.FlagDescList)
	InstanceCmd.Flags().StringP(config.FlagNameFind, config.FlagShortFind, "", config.FlagDescFind)
	InstanceCmd.Flags().BoolP(config.FlagNameImageDetails, config.FlagShortImageDetails, false, config.FlagDescImageDetails)
}

// setupInstanceContext initializes logging, creates AppContext, and registers subcommands.
func setupInstanceContext(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "help" {
		return nil
	}

	// Initialize logger
	if err := logger.SetLogger(); err != nil {
		return err
	}
	logger.InitLogger(logger.CmdLogger)

	// Create and inject AppContext
	ctx := cmd.Context()
	appCtx, err := app.NewAppContext(ctx, cmd)
	if err != nil {
		return fmt.Errorf("initializing app context: %w", err)
	}
	ctx = context.WithValue(ctx, "appCtx", appCtx)
	cmd.SetContext(ctx)

	// Register list/find as subcommands for newer usage
	if cmd.Name() == "instance" && len(cmd.Commands()) == 0 {
		cmd.AddCommand(
			newListCmd(appCtx),
			newFindCmd(appCtx),
		)
	}
	return nil
}

// executeInstanceCommand handles the old flag syntax for backward compatibility.
func executeInstanceCommand(cmd *cobra.Command, args []string) error {
	appCtx, err := getAppContext(cmd)
	if err != nil {
		return err
	}

	list, _ := cmd.Flags().GetBool(config.FlagNameList)
	find, _ := cmd.Flags().GetString(config.FlagNameFind)
	imageDetails, _ := cmd.Flags().GetBool(config.FlagNameImageDetails)

	switch {
	case list:
		fmt.Println("Listing instances in compartment:", appCtx.CompartmentName)
		return resources.ListInstances(appCtx)

	case find != "":
		fmt.Println("Finding instances with name pattern:", find)
		return resources.FindInstances(appCtx, find, imageDetails)

	default:
		return cmd.Help()
	}
}

// getAppContext retrieves the AppContext from the command's context or returns an error.
func getAppContext(cmd *cobra.Command) (*app.AppContext, error) {
	ctx := cmd.Context()
	appCtx, ok := ctx.Value("appCtx").(*app.AppContext)
	if !ok || appCtx == nil {
		return nil, fmt.Errorf("app context not found in command context")
	}
	return appCtx, nil
}
