package instance

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/pkg/flags"
	"github.com/rozdolsky33/ocloud/pkg/resources/compute"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
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
	flags.ListFlag.AddBoolFlag(InstanceCmd)
	flags.FindFlag.AddStringFlag(InstanceCmd)
	flags.ImageDetailsFlag.AddBoolFlag(InstanceCmd)
}

// setupInstanceContext initializes logging, creates AppContext, and registers subcommands.
func setupInstanceContext(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "help" {
		return nil
	}

	return initializeCommandContext(cmd)
}

// initializeCommandContext handles logger initialization, AppContext creation, and subcommand registration.
func initializeCommandContext(cmd *cobra.Command) error {
	// Initialize logger
	if err := logger.SetLogger(); err != nil {
		return err
	}
	logger.InitLogger(logger.CmdLogger)

	// Create AppContext
	ctx := cmd.Context()
	application, err := app.InitApp(ctx, cmd)
	if err != nil {
		return fmt.Errorf("initializing app: %w", err)
	}

	// Store AppContext in command's context for backward compatibility and exploration mode
	// This will be removed in future versions.
	ctx = context.WithValue(ctx, "appCtx", application)
	cmd.SetContext(ctx)

	// Register list/find as subcommands for newer usage
	if cmd.Name() == "instance" && len(cmd.Commands()) == 0 {
		cmd.AddCommand(
			newListCmd(application),
			newFindCmd(application),
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
	return doInstanceCommand(cmd, appCtx)
}

// doInstanceCommand handles the actual execution of instance commands based on flags.
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

// getAppContext retrieves the AppContext from the command's context or returns an error.
func getAppContext(cmd *cobra.Command) (*app.AppContext, error) {
	ctx := cmd.Context()
	if ctx == nil {
		return nil, fmt.Errorf("command context is nil")
	}
	appCtx, ok := ctx.Value("appCtx").(*app.AppContext)
	if !ok || appCtx == nil {
		return nil, fmt.Errorf("app context not found in command context")
	}
	return appCtx, nil
}
