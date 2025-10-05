package objectstorage

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	os "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// GetBuckets retrieves and displays a paginated list of object storage buckets in a given compartment.
// It uses a specified limit and page number for pagination and can output results as JSON based on the flag.
// Returns an error if bucket retrieval or output processing fails.
func GetBuckets(appCtx *app.ApplicationContext, limit int, page int, useJSON bool) error {
	ctx := context.Background()
	client, err := oci.NewObjectStorageClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating object storage client: %w", err)
	}

	bucketAdapter := os.NewAdapter(client)
	service := NewService(bucketAdapter, appCtx.Logger, appCtx.CompartmentID)

	buckets, total, next, err := service.FetchPaginatedBuckets(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing buckets: %w", err)
	}

	return PrintBucketsInfo(buckets, appCtx, &util.PaginationInfo{CurrentPage: page, TotalCount: total, NextPageToken: next, Limit: limit}, useJSON)
}
