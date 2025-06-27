package policy

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/identity/policy"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
Find Policies in the specified compartment that match the given pattern.

The search is performed using a fuzzy matching algorithm that searches across multiple fields:

Searchable Fields:
- Name: Policy name
- Description: Description of the policy
- Statement: Policy statements

The search pattern is automatically wrapped with wildcards, so partial matches are supported.
For example, searching for "admin" will match "administrators", "admin-policy", etc.

You can also search for specific tag values by using the tag key and value in your search pattern.
For example, "environment:production" will find policies with that specific tag.
`

// Examples for the find command
var findExamples = `
  # Find Policies with "admin" in their name
  ocloud identity policy find admin

  # Find Policies with "network" in their name and output in JSON format
  ocloud identity policy find network --json
`

// NewFindCmd creates a new command for finding policies by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find Policies by name pattern",
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

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running policy find command", "pattern", namePattern, "json", useJSON)
	return policy.FindPolicies(appCtx, namePattern, useJSON)
}
