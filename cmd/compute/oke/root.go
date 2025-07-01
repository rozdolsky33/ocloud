package oke

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewOKECmd creates a new command for OKE-related operations
func NewOKECmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "oke",
		Short:         "Manage OCI Kubernetes Engine (OKE)",
		Long:          "Manage Oracle Cloud Infrastructure Kubernetes Engine (OKE) clusters and node pools.\n\nThis command allows you to list all clusters in a compartment or find specific clusters by name pattern. For each cluster, you can view detailed information including Kubernetes version, endpoint, and associated node pools.",
		Example:       "  ocloud compute oke list\n  ocloud compute oke list --json\n  ocloud compute oke find myoke\n  ocloud compute oke find myoke --json",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
