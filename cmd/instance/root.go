package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/pkg/resources/compute"
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// Store the application context to avoid double initialization
var appContext *app.AppContext

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
	config.ListFlag.AddBoolFlag(InstanceCmd)
	config.FindFlag.AddStringFlag(InstanceCmd)
	config.ImageDetailsFlag.AddBoolFlag(InstanceCmd)
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
	var err error
	appContext, err = app.InitApp(ctx, cmd)
	if err != nil {
		return fmt.Errorf("initializing app: %w", err)
	}

	// Register list/find as subcommands for newer usage
	if cmd.Name() == "instance" && len(cmd.Commands()) == 0 {
		cmd.AddCommand(
			newListCmd(appContext),
			newFindCmd(appContext),
		)
	}
	return nil
}

// executeInstanceCommand handles the old flag syntax for backward compatibility.
func executeInstanceCommand(cmd *cobra.Command, args []string) error {
	// Use the already initialized AppContext
	if appContext == nil {
		return fmt.Errorf("app context not initialized")
	}
	return doInstanceCommand(cmd, appContext)
}

// doInstanceCommand handles the actual execution of instance commands based on config.
func doInstanceCommand(cmd *cobra.Command, appCtx *app.AppContext) error {
	list, _ := cmd.Flags().GetBool(config.FlagNameList)
	find, _ := cmd.Flags().GetString(config.FlagNameFind)
	imageDetails, _ := cmd.Flags().GetBool(config.FlagNameImageDetails)

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
