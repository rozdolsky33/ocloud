package objectstorage

import (
	"context"
	"fmt"
	"time"

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

	db := a.toDomainBucketFromBucket(resp.Bucket)
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
	var tc *time.Time
	if bucket.TimeCreated != nil {
		t := bucket.TimeCreated.Time
		tc = &t
	}
	attrs := mapping.BucketAttributes{
		Name:        bucket.Name,
		Namespace:   bucket.Namespace,
		TimeCreated: tc,
	}
	return mapping.NewDomainBucketFromAttrs(attrs)
}

func (a *Adapter) toDomainBucketFromBucket(bucket objectstorage.Bucket) domain.Bucket {
	var tc *time.Time
	if bucket.TimeCreated != nil {
		t := bucket.TimeCreated.Time
		tc = &t
	}
	attrs := mapping.BucketAttributes{
		Name:               bucket.Name,
		ID:                 bucket.Id,
		Namespace:          bucket.Namespace,
		TimeCreated:        tc,
		StorageTier:        string(bucket.StorageTier),
		PublicAccessType:   string(bucket.PublicAccessType),
		KmsKeyID:           bucket.KmsKeyId,
		Versioning:         string(bucket.Versioning),
		ReplicationEnabled: bucket.ReplicationEnabled,
		IsReadOnly:         bucket.IsReadOnly,
		ApproximateCount:   bucket.ApproximateCount,
		ApproximateSize:    bucket.ApproximateSize,
		FreeformTags:       bucket.FreeformTags,
		DefinedTags:        bucket.DefinedTags,
	}
	return mapping.NewDomainBucketFromAttrs(attrs)
}
