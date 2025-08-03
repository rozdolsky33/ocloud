package setup

import (
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/setup"
	"github.com/spf13/cobra"
)

// SetupMappingFile creates a new command for setting up or updating the tenancy mapping file.
func SetupMappingFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "setup",
		Short:         "Create tenancy mapping file or add a record",
		Long:          "",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSetupFileMappingCommand(cmd)
		},
	}

	return cmd
}

func RunSetupFileMappingCommand(cmd *cobra.Command) error {
	logger.LogWithLevel(logger.CmdLogger, 1, "Running setup command")
	return setup.SetupTenancyMapping()
}
