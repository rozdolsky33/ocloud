package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

func NewObjectStorageCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "object-storage",
		Aliases:       []string{"objectstorage", "os"},
		Short:         "Manage OCI Object Storage: list, get, and search",
		Long:          "Manage Oracle Cloud Infrastructure Object Storage: list, get, and search\",",
		Example:       "  ocloud storage object-storage list\n  ocloud storage object-storage list --json\n  ocloud storage object-storage get\n  ocloud storage object-storage get --json\n  ocloud storage object-storage search <value>\n  ocloud storage object-storage search <value> --json",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	return cmd
}
