package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	domain "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
)

// BucketAttributes is a generic, intermediate representation of a bucket's data.
// Each adapter is responsible for populating this struct.
type BucketAttributes struct {
	Name               *string
	ID                 *string
	Namespace          *string
	TimeCreated        *time.Time
	StorageTier        string
	PublicAccessType   string
	KmsKeyID           *string
	Versioning         string
	ReplicationEnabled *bool
	IsReadOnly         *bool
	ApproximateCount   *int64
	ApproximateSize    *int64
	FreeformTags       map[string]string
	DefinedTags        map[string]map[string]interface{}
}

func NewBucketAttributesFromOCIBucket(bucket objectstorage.Bucket) *BucketAttributes {
	var tc *time.Time
	if bucket.TimeCreated != nil {
		t := bucket.TimeCreated.Time
		tc = &t
	}
	return &BucketAttributes{
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
}

func NewBucketAttributesFromOCIBucketSummary(bucket objectstorage.BucketSummary) *BucketAttributes {
	var tc *time.Time
	if bucket.TimeCreated != nil {
		t := bucket.TimeCreated.Time
		tc = &t
	}
	return &BucketAttributes{
		Name:        bucket.Name,
		Namespace:   bucket.Namespace,
		TimeCreated: tc,
	}
}

// NewDomainBucketFromAttrs builds a domain.Bucket from provider-agnostic attributes.
func NewDomainBucketFromAttrs(attrs BucketAttributes) domain.Bucket {
	encryption := "Oracle-managed"
	if kmsKeyID := stringValue(attrs.KmsKeyID); kmsKeyID != "" {
		encryption = "Customer-managed (KMS)"
	}

	db := domain.Bucket{
		Name:               stringValue(attrs.Name),
		OCID:               stringValue(attrs.ID),
		Namespace:          stringValue(attrs.Namespace),
		StorageTier:        attrs.StorageTier,
		Visibility:         attrs.PublicAccessType,
		Versioning:         attrs.Versioning,
		Encryption:         encryption,
		ReplicationEnabled: boolValue(attrs.ReplicationEnabled),
		IsReadOnly:         boolValue(attrs.IsReadOnly),
		ApproximateCount:   intValueFromInt64(attrs.ApproximateCount),
		ApproximateSize:    int64Value(attrs.ApproximateSize),
		FreeformTags:       attrs.FreeformTags,
		DefinedTags:        attrs.DefinedTags,
	}

	if attrs.TimeCreated != nil {
		db.TimeCreated = *attrs.TimeCreated
	}

	return db
}

// Helper to dereference a *string safely.
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper to dereference a *bool safely.
func boolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Helper to convert *int64 to int safely.
func intValueFromInt64(v *int64) int {
	if v == nil {
		return 0
	}
	return int(*v)
}

// Helper to dereference *int64 safely.
func int64Value(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}
