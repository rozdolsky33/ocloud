package info

import (
	configurationFlags "github.com/rozdolsky33/ocloud/cmd/configuration/flags"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/info"
	"github.com/spf13/cobra"
)

// Long description for the map-file command
var mapFileLong = `
View the tenancy mapping information from the tenancy-map.yaml file.

This command displays information about the tenancy mappings defined in the tenancy-map.yaml file.
It shows details such as environment, tenancy, tenancy ID, realm, compartments, and regions.

Additional Information:
- Use --json (-j) to output the results in JSON format
- Use --realm (-r) to filter the mappings by realm (e.g., OC1, OC2, etc.)
- The command reads the tenancy-map.yaml file from the default location or from the path specified by the OCI_TENANCY_MAP_PATH environment variable
`

// Examples for the map-file command
var mapFileExamples = `
  # View the tenancy mapping information
  ocloud config info map-file

  # View the tenancy mapping information in JSON format
  ocloud config info map-file --json

  # Filter tenancy mappings by realm
  ocloud config info map-file --realm OC1

  # Filter tenancy mappings by realm and output in JSON format
  ocloud config info map-file --realm OC1 --json
`

// ViewMappingFile creates a new command for viewing the tenancy mapping file
func ViewMappingFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "map-file",
		Aliases:       []string{"mf", "tf"},
		Short:         "View tenancy mapping information",
		Long:          mapFileLong,
		Example:       mapFileExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunViewFileMappingCommand(cmd)
		},
	}

	// Add JSON flag
	flags.JSONFlag.Add(cmd)

	// Add realm filter flag
	configurationFlags.RealmFlag.Add(cmd)

	return cmd
}

// RunViewFileMappingCommand handles the execution of the map-file command
func RunViewFileMappingCommand(cmd *cobra.Command) error {
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	realm := flags.GetStringFlag(cmd, flags.FlagNameRealm, "")
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running map-file command", "json", useJSON, "realm", realm)
	err := info.ViewConfiguration(useJSON, realm)
	if err != nil {
		return err
	}
	logger.CmdLogger.V(logger.Info).Info("Map-file command completed.")
	return nil
}
