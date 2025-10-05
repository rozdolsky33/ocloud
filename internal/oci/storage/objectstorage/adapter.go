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

// GetBucketNameByOCID retrieves the OCID of a bucket by its name.
func (a *Adapter) GetBucketNameByOCID(ctx context.Context, compartmentID, bucketOCID string) (string, error) {
	nsResp, err := a.client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get namespace: %w", err)
	}

	namespace := *nsResp.Value
	var page *string

	for {
		listResp, err := a.client.ListBuckets(ctx, objectstorage.ListBucketsRequest{
			NamespaceName: &namespace,
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return "", fmt.Errorf("list buckets: %w", err)
		}

		for _, sum := range listResp.Items {
			if sum.Name == nil || *sum.Name == "" {
				continue
			}
			name := *sum.Name

			getResp, err := a.client.GetBucket(ctx, objectstorage.GetBucketRequest{
				NamespaceName: &namespace,
				BucketName:    &name,
			})
			if err != nil || getResp.Bucket.Id == nil {
				continue
			}
			if *getResp.Bucket.Id == bucketOCID {
				return name, nil
			}
		}

		if listResp.OpcNextPage == nil {
			break
		}
		page = listResp.OpcNextPage
	}

	return "", fmt.Errorf("bucket with OCID %q not found in compartment %q", bucketOCID, compartmentID)
}

// GetBucketByName retrieves a single bucket by its name (interprets input string as bucket name).
func (a *Adapter) GetBucketByName(ctx context.Context, bucketName string) (*domain.Bucket, error) {
	nsResp, err := a.client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}
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

	bkt := mapping.NewDomainBucketFromAttrs(*mapping.NewBucketAttributesFromOCIBucket(resp.Bucket))

	return &bkt, nil
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
