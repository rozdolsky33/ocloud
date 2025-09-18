package vcn

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	netvcn "github.com/rozdolsky33/ocloud/internal/services/network/vcn"
	"github.com/spf13/cobra"
)

// Long description for the get command
var getLong = `
Display a VCN summary by OCID.

This command fetches a Virtual Cloud Network (VCN) by its OCID and prints a concise summary
including name, state, compartment, CIDR blocks, DNS label/domain, DHCP options, creation time,
and tags. Use --json to get the raw JSON output instead of the formatted summary.`

// Examples for the get command
var getExamples = `
  # Show VCN summary
  ocloud network vcn get ocid1.vcn.oc1..aaaaaaaa...

  # Show VCN summary in JSON
  ocloud network vcn get ocid1.vcn.oc1..aaaaaaaa... --json
`

// NewGetCmd returns "vcn get" command.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get <vcn-ocid>",
		Short:         "Get VCN summary by OCID",
		Long:          getLong,
		Example:       getExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunGetCommand executes the get logic
func RunGetCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	vcnID := args[0]
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn get", "vcnID", vcnID, "json", useJSON)

	return netvcn.GetVCN(appCtx, vcnID, useJSON)
}
