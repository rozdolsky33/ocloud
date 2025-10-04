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
	Visibility         string
	Encryption         string
	Versioning         string
	ReplicationEnabled bool
	IsReadOnly         bool
	ApproximateSize    int64
	ApproximateCount   int
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
