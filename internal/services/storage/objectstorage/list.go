package objectstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	os "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListBuckets retrieves and lists all buckets in the specified compartment. It supports both TUI and JSON output formatting.
func ListBuckets(appCtx *app.ApplicationContext, useJSON bool) error {
	ctx := context.Background()
	client, err := oci.NewObjectStorageClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating object storage client: %w", err)
	}

	bucketAdapter := os.NewAdapter(client)
	service := NewService(bucketAdapter, appCtx.Logger, appCtx.CompartmentID)
	buckets, err := service.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("listing buckets: %w", err)
	}

	//TUI
	model := os.NewBucketListModel(buckets)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting bucket: %w", err)
	}

	name, err := service.osRepo.GetBucketNameByOCID(ctx, appCtx.CompartmentID, id)
	bucket, err := service.osRepo.GetBucketByName(ctx, appCtx.CompartmentID, name)
	if err != nil {
		return fmt.Errorf("getting bucket: %w", err)
	}

	return PrintBucketInfo(bucket, appCtx, useJSON)
}
