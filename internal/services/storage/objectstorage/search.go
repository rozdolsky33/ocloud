package objectstorage

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociobj "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
)

func SearchBuckets(appCtx *app.ApplicationContext, pattern string, useJSON bool) error {
	ctx := context.Background()
	client, err := oci.NewObjectStorageClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating object storage client: %w", err)
	}

	adapter := ociobj.NewAdapter(client)
	svc := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	buckets, err := svc.FuzzySearch(ctx, pattern)
	if err != nil {
		return fmt.Errorf("searching buckets: %w", err)
	}
	if err := PrintBucketsInfo(buckets, appCtx, nil, useJSON); err != nil {
		return fmt.Errorf("printing buckets: %w", err)
	}
	logger.LogWithLevel(logger.CmdLogger, logger.Info, "Found matching buckets", "search", pattern, "matched", len(buckets))
	return nil
}
