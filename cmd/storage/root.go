package storage

import (
	"github.com/rozdolsky33/ocloud/cmd/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

func NewStorageCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "storage",
		Aliases:       []string{"stg"},
		Short:         "Manage OCI Storage Resources",
		Long:          "Manage Oracle Cloud Infrastructure Storage Resources: list, get, and search by name or pattern.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(objectstorage.NewObjectStorageCmd(appCtx))
	return cmd
}
