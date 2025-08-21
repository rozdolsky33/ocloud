package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	autonomousdboci "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
	"github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
Find Autonomous Databases in the specified compartment that match the given pattern.

The search is performed using a fuzzy matching algorithm that searches across multiple fields:

Searchable Fields:
- Name: Database name
- DisplayName: Display name of the database
- DbName: Database name

The search pattern is automatically wrapped with wildcards, so partial matches are supported.
For example, searching for "prod" will match "production", etc.
`

// Examples for the find command
var findExamples = `
  # Find Autonomous Databases with "prod" in their name
  ocloud database autonomous find prod

  # Find Autonomous Databases with "test" in their name and output in JSON format
  ocloud database autonomous find test --json
`

// NewFindCmd creates a new command for finding compartments by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find Database by name pattern",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, 1, "Running find command", "pattern", namePattern, "json", useJSON)

	repo, err := autonomousdboci.NewAdapter(appCtx.Provider, appCtx.CompartmentID)
	if err != nil {
		return err
	}
	service := autonomousdb.NewService(repo, appCtx)
	databases, err := service.Find(cmd.Context(), namePattern)
	if err != nil {
		return err
	}
	// Convert service type to domain type for output
	domainDbs := make([]domain.AutonomousDatabase, 0, len(databases))
	for _, db := range databases {
		domainDbs = append(domainDbs, domain.AutonomousDatabase(db))
	}
	return autonomousdb.PrintAutonomousDbInfo(domainDbs, appCtx, nil, useJSON)
}
