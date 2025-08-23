package setup

import (
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/setup"
	"github.com/spf13/cobra"
)

// Long description for the setup command
var setupLong = `
Create or update the tenancy mapping file used by ocloud CLI.

This command guides you through an interactive process to create a new tenancy mapping file
or add records to an existing one. The mapping file allows ocloud to associate tenancy names
with their OCIDs and other metadata such as compartments and regions.

The tenancy mapping file is stored at ~/.oci/.ocloud/tenancy-map.yaml by default, but this
location can be overridden using the OCI_TENANCY_MAP_PATH environment variable.

Each record in the mapping file includes:
- Environment: A descriptive name for the environment (e.g., Prod, Dev, Test)
- Tenancy Name: The name of the tenancy
- Tenancy OCID: The Oracle Cloud ID of the tenancy
- Realm: The OCI realm (e.g., OC1, OC2)
- Compartments: A list of compartments in the tenancy
- Regions: A list of regions used by the tenancy
`

// Examples for the setup command
var setupExamples = `
  # Create or update the tenancy mapping file
  ocloud config setup

  # After running the command, you'll be guided through an interactive process
  # to enter information about your tenancy environments
`

// SetupMappingFile creates a new command for setting up or updating the tenancy mapping file.
func SetupMappingFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "setup",
		Short:         "Create tenancy mapping file or add a record",
		Long:          setupLong,
		Example:       setupExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSetupFileMappingCommand(cmd)
		},
	}

	return cmd
}

// RunSetupFileMappingCommand handles the execution of the setup command
func RunSetupFileMappingCommand(cmd *cobra.Command) error {
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running setup command")
	err := setup.SetupTenancyMapping()
	if err != nil {
		return err
	}
	logger.CmdLogger.V(logger.Info).Info("Setup command completed.")
	return nil
}
