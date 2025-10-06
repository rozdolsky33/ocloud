package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

var searchLong = `
Search for instances in the specified compartment that match the given pattern.

The search uses a fuzzy, prefix, and substring matching algorithm across many indexed fields.
You can search using any of the following fields (partial matches are supported):

Searchable fields:
- Name: Display name of the instance
- Hostname: Instance hostname
- ImageName: Name of the image used by the instance
- ImageOS: Operating system of the image
- Shape: Instance shape
- PrimaryIP: Primary private IP address
- OCID: Instance OCID
- VcnName: Name of the VCN the instance is attached to
- SubnetName: Name of the subnet the instance is attached to
- FD: Fault domain
- AD: Availability domain
- TagsKV: All tags in key=value form, flattened
- TagsVal: Only tag values (e.g., "8.10")

The search pattern is case-insensitive. For very specific inputs (like full OCID, IP, or exact hostname),
the search first tries exact and substring matches; otherwise it falls back to broader fuzzy search.
`

var searchExamples = `
  # Search by display name (substring)
  ocloud compute instance search web

  # Search by hostname (substring)
  ocloud compute instance search host123

  # Search by image name or OS
  ocloud compute instance search oracle
  ocloud compute instance search "Oracle-Linux"

  # Search by shape
  ocloud compute instance search VM.Standard3

  # Search by primary IP (exact or partial)
  ocloud compute instance search 10.0.1.15
  ocloud compute instance search 10.0.1.

  # Search by OCID (exact)
  ocloud compute instance search ocid1.instance.oc1..aaaa...

  # Search by VCN or Subnet names
  ocloud compute instance search my-vcn
  ocloud compute instance search app-subnet

  # Search by Availability/Fault Domain
  ocloud compute instance search AD-1
  ocloud compute instance search FD-2

  # Search by tag value only (TagsVal)
  ocloud compute instance search 8.10

  # Show more details in the output
  ocloud compute instance search api --all

  # Output in JSON format
  ocloud compute instance search server --json
`

// NewSearchCmd creates a new command for finding instances by name pattern
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Fuzzy Search for Instances",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}

	instaceFlags.AllInfoFlag.Add(cmd)

	return cmd
}

// RunSearchCommand handles the execution of the find command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	showDetails := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running instance search command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return instance.SearchInstances(appCtx, namePattern, useJSON, showDetails)
}
