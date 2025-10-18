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
		Short:         "Explore OCI Storage Resources",
		Long:          "Explore Oracle Cloud Infrastructure Storage Resources: object storage, block storage, file storage, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(objectstorage.NewObjectStorageCmd(appCtx))
	return cmd
}
