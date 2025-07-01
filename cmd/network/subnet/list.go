package subnet

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/network/subnet"
	"github.com/spf13/cobra"
)

// Long description for the list command
var listLong = `
List all Subnets in the specified tenancy or compartment.

This command displays information about all subnets in the current compartment,
including their names, CIDR blocks, and whether they allow public IP addresses.
By default, it shows basic subnet information in a tabular format.

Additional Information:
- Use --json (-j) to output the results in JSON format
- Use --limit (-m) to control the number of results per page
- Use --page (-p) to navigate between pages of results
- Use --sort (-s) to sort results by name or CIDR
`

// Examples for the list command
var listExamples = `
  # List all subnets in the current compartment
  ocloud network subnet list

  # List all subnets and output in JSON format
  ocloud network subnet list --json

  # List subnets with pagination (10 per page, page 2)
  ocloud network subnet list --limit 10 --page 2

  # List subnets sorted by name
  ocloud network subnet list --sort name

  # List subnets sorted by CIDR block
  ocloud network subnet list --sort cidr
`

// NewListCmd creates a new "list" command for listing all subnets in a specified tenancy or compartment.
// Accepts an ApplicationContext for accessing configuration and dependencies.
// Adds pagination flags for controlling the number of results returned and the page to retrieve.
// Executes the RunListCommand function when invoked.
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Subnets in the specified tenancy or compartment",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}
	// Add flags specific to the list command
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)
	paginationFlags.SortFlag.Add(cmd)

	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Get pagination parameters
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	sortBy := flags.GetStringFlag(cmd, flags.FlagNameSort, "")

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running subnet list command in", "compartment", appCtx.CompartmentName, "json", useJSON, "sort", sortBy)
	return subnet.ListSubnets(appCtx, useJSON, limit, page, sortBy)
}
