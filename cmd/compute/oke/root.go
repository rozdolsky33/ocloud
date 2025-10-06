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
		Long:          "Manage Oracle Cloud Infrastructure Kubernetes Engine (OKE) clusters and node pools.\nThis command allows you to list all clusters in a compartment or search specific clusters by search pattern. For each cluster, you can view detailed information including Kubernetes version, endpoint, and associated node pools.",
		Example:       "  ocloud compute oke list\n  ocloud compute oke list --json\n  ocloud compute oke get\n  ocloud compute oke get --json\n  ocloud compute oke search myoke\n  ocloud compute oke search myoke --json",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))

	return cmd
}
