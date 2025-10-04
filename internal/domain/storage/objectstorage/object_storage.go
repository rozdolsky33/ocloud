package objectstorage

import (
	"context"
	"time"
)

type Bucket struct {
	Name               string
	OCID               string
	Namespace          string
	StorageTier        string
	Visibility         string // e.g., NoPublicAccess, ObjectRead
	Encryption         string // e.g., Oracle-managed, Customer-managed
	Versioning         string // e.g., Enabled, Suspended
	ReplicationEnabled bool
	IsReadOnly         bool
	ApproximateSize    int64 // bytes
	ApproximateCount   int   // number of objects
	TimeCreated        time.Time
	FreeformTags       map[string]string
	DefinedTags        map[string]map[string]interface{}
}

type Object struct {
	Name         string
	Size         int64
	StorageTier  string
	ContentType  string
	LastModified time.Time
	ETag         string
}

type ObjectStorageRepository interface {
	GetBucket(ctx context.Context, ocid string) (*Bucket, error)
	ListBuckets(ctx context.Context, compartmentID string) ([]Bucket, error)
	//GetObject(ctx context.Context, bucketName, objectName string) (*Object, error)
	//ListObjects(ctx context.Context, bucketName string) ([]Object, error)
}
