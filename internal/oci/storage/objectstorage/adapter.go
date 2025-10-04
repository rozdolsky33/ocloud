package objectstorage

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	domain "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
)

type Adapter struct {
	client objectstorage.ObjectStorageClient
}

func (a *Adapter) GetBucket(ctx context.Context, name string) (*domain.Bucket, error) {
	nsResp, err := a.client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	resp, err := a.client.GetBucket(ctx, objectstorage.GetBucketRequest{
		NamespaceName: nsResp.Value,
		BucketName:    &name,
	})
	if err != nil {
		return nil, err
	}

	db := a.toDomainBucketFromBucket(resp.Bucket)
	return &db, nil
}

func NewAdapter(client objectstorage.ObjectStorageClient) *Adapter {
	return &Adapter{client: client}
}

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
			allBuckets = append(allBuckets, a.toDomainBucketFromSummary(item))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return allBuckets, nil
}

func (a *Adapter) toDomainBucketFromSummary(bucket objectstorage.BucketSummary) domain.Bucket {
	return domain.Bucket{
		Name:        *bucket.Name,
		Namespace:   *bucket.Namespace,
		TimeCreated: bucket.TimeCreated.Time,
	}
}

func (a *Adapter) toDomainBucketFromBucket(bucket objectstorage.Bucket) domain.Bucket {
	return domain.Bucket{
		Name:        *bucket.Name,
		OCID:        *bucket.Id,
		Namespace:   *bucket.Namespace,
		TimeCreated: bucket.TimeCreated.Time,
	}
}
