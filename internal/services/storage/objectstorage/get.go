package objectstorage

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	os "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
)

func GetBuckets(appCtx *app.ApplicationContext, limit int, page int, useJSON bool) error {
	ctx := context.Background()
	client, err := oci.NewObjectStorageClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating object storage client: %w", err)
	}

	bucketAdapter := os.NewAdapter(client)
	service := NewService(bucketAdapter, appCtx.Logger, appCtx.CompartmentID)

	buckets, err := service.osRepo.ListBuckets(ctx, service.CompartmentID)
	if err != nil {
		return err
	}

	return PrintBucketsInfo(buckets, appCtx, nil, useJSON)
}
