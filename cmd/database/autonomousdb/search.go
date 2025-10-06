package autonomousdb

import (
	databaseFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	"github.com/spf13/cobra"
)

// Long description for the search command
var findLong = `
FuzzySearch Autonomous Databases in the specified compartment that match the given pattern.

This command searches across a broad set of Autonomous Database attributes. It matches on:
- Identity & lifecycle: Name, OCID, State, DB version, Workload, License model, Compute model
- Capacity: OCPU/ECPU count, CPU core count, Storage size (GB/TB)
- Networking: VCN ID/Name, Subnet ID/Name, Private endpoint, Private endpoint IP/label,
  Whitelisted IPs, NSG IDs/Names
- Tags: Flattened key:value pairs (TagsKV) and tag values (TagsVal)

The search is case-insensitive and uses a fuzzy/prefix/wildcard strategy, similar to instance and OKE search.

Additional Information:
- Use --json (-j) to output results in JSON format
- Works with partial fragments (e.g., OCID parts, hostnames, IPs, tag values)
`

// Examples for the search command
var findExamples = `
  # Fuzzy search ADBs with names containing "prod"
  ocloud database autonomous search prod

  # Output results in JSON
  ocloud database autonomous search dev --json

  # Search by OCID fragment or exact OCID
  ocloud database autonomous search ocid1.autonomousdatabase
  ocloud database autonomous search ocid1.autonomousdatabase.oc1..exampleexactocid

  # Search by DB version, workload, license, or compute model
  ocloud database autonomous search 19c
  ocloud database autonomous search OLTP
  ocloud database autonomous search LICENSE_INCLUDED
  ocloud database autonomous search ECPU

  # Search by networking attributes
  ocloud database autonomous search my-vcn
  ocloud database autonomous search app-subnet
  ocloud database autonomous search adb-priv.endpoint.oraclecloud.com
  ocloud database autonomous search 10.0.1.25

  # Search by NSG name or whitelisted IP
  ocloud database autonomous search nsg-apps
  ocloud database autonomous search 203.0.113.42

  # Search by tags (key:value or value only)
  ocloud database autonomous search team:platform
  ocloud database autonomous search prod
`

// NewFindCmd creates a new command for finding compartments by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for Autonomous Databases",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFindCommand(cmd, args, appCtx)
		},
	}
	databaseFlags.AllInfoFlag.Add(cmd)
	return cmd
}

// runFindCommand handles the execution of the find command
func runFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running search command", "searchPattern", namePattern, "json", useJSON)
	showAll := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	return autonomousdb.SearchAutonomousDatabases(appCtx, namePattern, useJSON, showAll)
}
