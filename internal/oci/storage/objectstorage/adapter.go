package objectstorage

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	domain "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/mapping"
)

type Adapter struct {
	client objectstorage.ObjectStorageClient
}

// NewAdapter builds a new adapter.
func NewAdapter(client objectstorage.ObjectStorageClient) *Adapter {
	return &Adapter{client: client}
}

// GetBucket retrieves a single bucket by its name (interprets input string as bucket name).
func (a *Adapter) GetBucket(ctx context.Context, ocid string) (*domain.Bucket, error) {
	nsResp, err := a.client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	bucketName := ocid
	resp, err := a.client.GetBucket(ctx, objectstorage.GetBucketRequest{
		NamespaceName: nsResp.Value,
		BucketName:    &bucketName,
		Fields: []objectstorage.GetBucketFieldsEnum{
			objectstorage.GetBucketFieldsApproximatecount,
			objectstorage.GetBucketFieldsApproximatesize,
		},
	})
	if err != nil {
		return nil, err
	}

	db := mapping.NewDomainBucketFromAttrs(*mapping.NewBucketAttributesFromOCIBucket(resp.Bucket))
	return &db, nil
}

// ListBuckets retrieves all buckets in a given compartment.
func (a *Adapter) ListBuckets(ctx context.Context, ocid string) (buckets []domain.Bucket, err error) {
	var allBuckets []domain.Bucket
	var page *string
	nsResp, err := a.client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	for {
		resp, err := a.client.ListBuckets(ctx, objectstorage.ListBucketsRequest{
			NamespaceName: nsResp.Value,
			CompartmentId: &ocid,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing buckets: %w", err)
		}
		for _, item := range resp.Items {
			allBuckets = append(allBuckets, mapping.NewDomainBucketFromAttrs(*mapping.NewBucketAttributesFromOCIBucketSummary(item)))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return allBuckets, nil
}
