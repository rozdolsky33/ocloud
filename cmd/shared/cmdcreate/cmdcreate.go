package cmdcreate

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/cmd/compute"
	"github.com/rozdolsky33/ocloud/cmd/configuration"
	"github.com/rozdolsky33/ocloud/cmd/database"
	"github.com/rozdolsky33/ocloud/cmd/identity"
	"github.com/rozdolsky33/ocloud/cmd/network"
	"github.com/rozdolsky33/ocloud/cmd/version"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
)

// CreateRootCmd creates a root command with or without application context
// If appCtx is nil, only commands that don't need context are added
// If appCtx is not nil, all commands are added
func CreateRootCmd(appCtx *app.ApplicationContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "ocloud",
		Short:        "Interact with Oracle Cloud Infrastructure",
		Long:         "",
		SilenceUsage: true,
	}

	// Initialize global flags
	flags.AddGlobalFlags(rootCmd)

	// Add commands that don't need context
	rootCmd.AddCommand(version.NewVersionCommand())
	version.AddVersionFlag(rootCmd)
	rootCmd.AddCommand(configuration.NewConfigCmd())

	// If appCtx is not nil, add commands that need context
	if appCtx != nil {
		rootCmd.AddCommand(compute.NewComputeCmd(appCtx))
		rootCmd.AddCommand(identity.NewIdentityCmd(appCtx))
		rootCmd.AddCommand(database.NewDatabaseCmd(appCtx))
		rootCmd.AddCommand(network.NewNetworkCmd(appCtx))
	}

	return rootCmd
}

// CreateRootCmdWithoutContext creates a root command without application context
// This is used for commands that don't need a full context
func CreateRootCmdWithoutContext() *cobra.Command {
	rootCmd := CreateRootCmd(nil)

	addPlaceholderCommands(rootCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	return rootCmd
}

// addPlaceholderCommands adds placeholder commands that will be displayed in help
// but will show a message about needing to initialize if they're actually run
func addPlaceholderCommands(rootCmd *cobra.Command) {
	commandTypes := []struct {
		use   string
		short string
	}{
		{"compute", "Manage OCI compute services"},
		{"identity", "Manage OCI identity services"},
		{"database", "Manage OCI Database services"},
		{"network", "Manage OCI networking services"},
	}

	for _, cmdType := range commandTypes {
		cmd := &cobra.Command{
			Use:   cmdType.use,
			Short: cmdType.short,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("this command requires application initialization")
			},
		}
		rootCmd.AddCommand(cmd)
	}
}
