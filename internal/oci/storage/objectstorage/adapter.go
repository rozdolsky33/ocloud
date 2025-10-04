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

func (a *Adapter) GetBucket(ctx context.Context, ocid string) (*domain.Bucket, error) {
	nsResp, err := a.client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
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
	db := domain.Bucket{}
	if bucket.Name != nil {
		db.Name = *bucket.Name
	}
	if bucket.Namespace != nil {
		db.Namespace = *bucket.Namespace
	}
	if bucket.TimeCreated != nil {
		db.TimeCreated = bucket.TimeCreated.Time
	}
	return db
}

func (a *Adapter) toDomainBucketFromBucket(bucket objectstorage.Bucket) domain.Bucket {
	db := domain.Bucket{}
	if bucket.Name != nil {
		db.Name = *bucket.Name
	}
	if bucket.Id != nil {
		db.OCID = *bucket.Id
	}
	if bucket.Namespace != nil {
		db.Namespace = *bucket.Namespace
	}
	if bucket.TimeCreated != nil {
		db.TimeCreated = bucket.TimeCreated.Time
	}
	if bucket.StorageTier != "" {
		db.StorageTier = string(bucket.StorageTier)
	}
	if bucket.PublicAccessType != "" {
		db.Visibility = string(bucket.PublicAccessType)
	}
	// Encryption
	if bucket.KmsKeyId != nil && *bucket.KmsKeyId != "" {
		db.Encryption = "Customer-managed (KMS)"
	} else {
		db.Encryption = "Oracle-managed"
	}
	// Versioning
	if bucket.Versioning != "" {
		db.Versioning = string(bucket.Versioning)
	}
	// Flags
	if bucket.ReplicationEnabled != nil {
		db.ReplicationEnabled = *bucket.ReplicationEnabled
	}
	if bucket.IsReadOnly != nil {
		db.IsReadOnly = *bucket.IsReadOnly
	}
	// Approximate metrics
	if bucket.ApproximateCount != nil {
		db.ApproximateCount = int(*bucket.ApproximateCount)
	}
	if bucket.ApproximateSize != nil {
		db.ApproximateSize = *bucket.ApproximateSize
	}
	// Tags
	if bucket.FreeformTags != nil {
		db.FreeformTags = bucket.FreeformTags
	}
	if bucket.DefinedTags != nil {
		db.DefinedTags = bucket.DefinedTags
	}
	return db
}
