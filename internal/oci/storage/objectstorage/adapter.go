package objectstorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	nsResp, err := a.client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{
		CompartmentId: &compartmentID,
	})

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
func (a *Adapter) GetBucketByName(ctx context.Context, compartmentID, bucketName string) (*domain.Bucket, error) {
	nsResp, err := a.client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{
		CompartmentId: &compartmentID,
	})
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
	nsResp, err := a.client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{
		CompartmentId: &ocid,
	})
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

// GetNamespace retrieves the object storage namespace for a compartment.
func (a *Adapter) GetNamespace(ctx context.Context, compartmentID string) (string, error) {
	nsResp, err := a.client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{
		CompartmentId: &compartmentID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get namespace: %w", err)
	}
	return *nsResp.Value, nil
}

// ListObjects retrieves all objects in a bucket.
func (a *Adapter) ListObjects(ctx context.Context, namespace, bucketName string) ([]domain.Object, error) {
	var allObjects []domain.Object
	var nextStart *string
	// Request additional fields: size, etag, timeModified, storageTier
	fields := "name,size,etag,timeCreated,timeModified,storageTier"

	for {
		resp, err := a.client.ListObjects(ctx, objectstorage.ListObjectsRequest{
			NamespaceName: &namespace,
			BucketName:    &bucketName,
			Start:         nextStart,
			Fields:        &fields,
		})
		if err != nil {
			return nil, fmt.Errorf("listing objects: %w", err)
		}

		for _, item := range resp.Objects {
			obj := mapping.NewDomainObjectFromAttrs(*mapping.NewObjectAttributesFromOCIObjectSummary(item, bucketName, namespace))
			allObjects = append(allObjects, obj)
		}

		if resp.NextStartWith == nil || *resp.NextStartWith == "" {
			break
		}
		nextStart = resp.NextStartWith
	}

	return allObjects, nil
}

// GetObjectHead retrieves object metadata using HeadObject.
func (a *Adapter) GetObjectHead(ctx context.Context, namespace, bucketName, objectName string) (*domain.Object, error) {
	resp, err := a.client.HeadObject(ctx, objectstorage.HeadObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
	})
	if err != nil {
		return nil, fmt.Errorf("head object: %w", err)
	}

	attrs := mapping.NewObjectAttributesFromOCIHeadObject(resp, bucketName, namespace)
	// Name is not in HeadObject response, set it manually
	attrs.Name = &objectName
	// Size is in ContentLength
	if resp.ContentLength != nil {
		attrs.Size = resp.ContentLength
	}

	obj := mapping.NewDomainObjectFromAttrs(*attrs)
	return &obj, nil
}

// DownloadObject downloads an object to the specified destination path with progress reporting.
func (a *Adapter) DownloadObject(ctx context.Context, namespace, bucketName, objectName, destPath string, progressFn func(domain.TransferProgress)) error {
	resp, err := a.client.GetObject(ctx, objectstorage.GetObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
	})
	if err != nil {
		return fmt.Errorf("get object: %w", err)
	}
	defer resp.Content.Close()

	// Use only the base name from the object name for the destination file
	baseName := filepath.Base(objectName)
	fullPath := filepath.Join(destPath, baseName)

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", fullPath, err)
	}
	defer file.Close()

	totalSize := int64(0)
	if resp.ContentLength != nil {
		totalSize = *resp.ContentLength
	}

	// If we have a progress callback and know the total size, use progress writer
	if progressFn != nil && totalSize > 0 {
		pw := &progressWriter{
			total:      totalSize,
			downloaded: 0,
			onProgress: func(downloaded int64) {
				progressFn(domain.TransferProgress{
					BytesTransferred: downloaded,
					TotalBytes:       totalSize,
					PartNumber:       1,
					TotalParts:       1,
				})
			},
		}
		_, err = io.Copy(file, io.TeeReader(resp.Content, pw))
	} else {
		_, err = io.Copy(file, resp.Content)
	}

	if err != nil {
		return fmt.Errorf("write file %s: %w", fullPath, err)
	}

	return nil
}

// progressWriter tracks download progress.
type progressWriter struct {
	total      int64
	downloaded int64
	onProgress func(int64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.downloaded += int64(n)
	if pw.onProgress != nil {
		pw.onProgress(pw.downloaded)
	}
	return n, nil
}

// Multipart upload constants
const (
	// MinPartSize is the minimum part size for multipart upload (10 MiB)
	MinPartSize = 10 * 1024 * 1024
	// MaxPartSize is the maximum part size (50 MiB for better parallelism)
	MaxPartSize = 50 * 1024 * 1024
	// MultipartThreshold - files larger than this use multipart upload (10 MiB)
	MultipartThreshold = 10 * 1024 * 1024
)

// UploadObject uploads a file to object storage using multipart upload for large files.
func (a *Adapter) UploadObject(ctx context.Context, namespace, bucketName, objectName, filePath string, progressFn func(domain.TransferProgress)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	fileSize := stat.Size()

	// Use simple upload for small files
	if fileSize <= MultipartThreshold {
		return a.simpleUpload(ctx, namespace, bucketName, objectName, file, fileSize, progressFn)
	}

	// Use multipart upload for large files
	return a.multipartUpload(ctx, namespace, bucketName, objectName, file, fileSize, progressFn)
}

// simpleUpload performs a single PUT request for small files.
func (a *Adapter) simpleUpload(ctx context.Context, namespace, bucketName, objectName string, file *os.File, fileSize int64, progressFn func(domain.TransferProgress)) error {
	contentType := detectContentType(file.Name())

	_, err := a.client.PutObject(ctx, objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
		ContentLength: &fileSize,
		PutObjectBody: file,
		ContentType:   &contentType,
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	if progressFn != nil {
		progressFn(domain.TransferProgress{
			BytesTransferred: fileSize,
			TotalBytes:       fileSize,
			PartNumber:       1,
			TotalParts:       1,
		})
	}

	return nil
}

// multipartUpload performs a multipart upload for large files.
func (a *Adapter) multipartUpload(ctx context.Context, namespace, bucketName, objectName string, file *os.File, fileSize int64, progressFn func(domain.TransferProgress)) error {
	// Calculate part size and count
	partSize := int64(MaxPartSize)
	partCount := int((fileSize + partSize - 1) / partSize)

	contentType := detectContentType(file.Name())

	// Create multipart upload
	createResp, err := a.client.CreateMultipartUpload(ctx, objectstorage.CreateMultipartUploadRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		CreateMultipartUploadDetails: objectstorage.CreateMultipartUploadDetails{
			Object:      &objectName,
			ContentType: &contentType,
		},
	})
	if err != nil {
		return fmt.Errorf("create multipart upload: %w", err)
	}

	uploadID := *createResp.UploadId

	// Upload parts
	var committedParts []objectstorage.CommitMultipartUploadPartDetails
	var bytesUploaded int64

	for partNum := 1; partNum <= partCount; partNum++ {
		offset := int64(partNum-1) * partSize
		size := partSize
		if offset+size > fileSize {
			size = fileSize - offset
		}

		// Create a section reader for this part wrapped as ReadCloser
		sectionReader := io.NewSectionReader(file, offset, size)
		partReader := io.NopCloser(sectionReader)

		resp, err := a.client.UploadPart(ctx, objectstorage.UploadPartRequest{
			NamespaceName:  &namespace,
			BucketName:     &bucketName,
			ObjectName:     &objectName,
			UploadId:       &uploadID,
			UploadPartNum:  &partNum,
			UploadPartBody: partReader,
			ContentLength:  &size,
		})
		if err != nil {
			// Abort on failure
			_, _ = a.client.AbortMultipartUpload(ctx, objectstorage.AbortMultipartUploadRequest{
				NamespaceName: &namespace,
				BucketName:    &bucketName,
				ObjectName:    &objectName,
				UploadId:      &uploadID,
			})
			return fmt.Errorf("upload part %d: %w", partNum, err)
		}

		committedParts = append(committedParts, objectstorage.CommitMultipartUploadPartDetails{
			PartNum: &partNum,
			Etag:    resp.ETag,
		})

		bytesUploaded += size
		if progressFn != nil {
			progressFn(domain.TransferProgress{
				BytesTransferred: bytesUploaded,
				TotalBytes:       fileSize,
				PartNumber:       partNum,
				TotalParts:       partCount,
			})
		}
	}

	// Commit multipart upload
	_, err = a.client.CommitMultipartUpload(ctx, objectstorage.CommitMultipartUploadRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
		UploadId:      &uploadID,
		CommitMultipartUploadDetails: objectstorage.CommitMultipartUploadDetails{
			PartsToCommit: committedParts,
		},
	})
	if err != nil {
		return fmt.Errorf("commit multipart upload: %w", err)
	}

	return nil
}

// detectContentType returns the MIME type based on file extension.
func detectContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		return "application/json"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".gz", ".gzip":
		return "application/gzip"
	case ".tar":
		return "application/x-tar"
	case ".xml":
		return "application/xml"
	case ".yaml", ".yml":
		return "application/x-yaml"
	default:
		return "application/octet-stream"
	}
}
