package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

func NewObjectStorageCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "object-storage",
		Aliases:       []string{"objectstorage", "os"},
		Short:         "Manage OCI Object Storage",
		Long:          "Manage Oracle Cloud Infrastructure Object Storage.\nThis command allows you to list all object storage in a compartment or find specific by name pattern. For each object storage, you can view detailed information.",
		Example:       "  ocloud storage object-storage list\n  ocloud storage object-storage list --json\n  ocloud storage object-storage get\n  ocloud storage object-storage get --json\n  ocloud storage object-storage find mybck\n  ocloud storage object-storage find mybkt --json",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	return cmd
}
