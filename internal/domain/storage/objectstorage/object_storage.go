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
	ContentMD5   string
	ETag         string
	LastModified time.Time
	// For URL generation
	BucketName string
	Namespace  string
}

// TransferProgress reports transfer progress (upload or download) to a callback function.
type TransferProgress struct {
	BytesTransferred int64
	TotalBytes       int64
	PartNumber       int
	TotalParts       int
}

// ObjectStorageRepository defines the port for interacting with object storage.
type ObjectStorageRepository interface {
	GetBucketNameByOCID(ctx context.Context, compartmentID, bucketOCID string) (string, error)
	GetBucketByName(ctx context.Context, compartmentID, name string) (*Bucket, error)
	ListBuckets(ctx context.Context, compartmentID string) ([]Bucket, error)
	GetNamespace(ctx context.Context, compartmentID string) (string, error)
	ListObjects(ctx context.Context, namespace, bucketName string) ([]Object, error)
	GetObjectHead(ctx context.Context, namespace, bucketName, objectName string) (*Object, error)
	DownloadObject(ctx context.Context, namespace, bucketName, objectName, destPath string, progressFn func(TransferProgress)) error
	UploadObject(ctx context.Context, namespace, bucketName, objectName, filePath string, progressFn func(TransferProgress)) error
}
