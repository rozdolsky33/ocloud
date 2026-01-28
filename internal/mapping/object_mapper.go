package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	domain "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
)

// ObjectAttributes is a generic, intermediate representation of an object's data.
type ObjectAttributes struct {
	Name         *string
	Size         *int64
	StorageTier  string
	ContentType  *string
	ContentMD5   *string
	ETag         *string
	LastModified *time.Time
	BucketName   string
	Namespace    string
}

// NewObjectAttributesFromOCIObjectSummary creates ObjectAttributes from an OCI ListObjects response item.
func NewObjectAttributesFromOCIObjectSummary(obj objectstorage.ObjectSummary, bucketName, namespace string) *ObjectAttributes {
	var lm *time.Time
	if obj.TimeModified != nil {
		t := obj.TimeModified.Time
		lm = &t
	}
	return &ObjectAttributes{
		Name:         obj.Name,
		Size:         obj.Size,
		StorageTier:  string(obj.StorageTier),
		ETag:         obj.Etag,
		LastModified: lm,
		BucketName:   bucketName,
		Namespace:    namespace,
	}
}

// NewObjectAttributesFromOCIHeadObject creates ObjectAttributes from an OCI HeadObject response.
func NewObjectAttributesFromOCIHeadObject(resp objectstorage.HeadObjectResponse, bucketName, namespace string) *ObjectAttributes {
	var lm *time.Time
	if resp.LastModified != nil {
		t := resp.LastModified.Time
		lm = &t
	}
	return &ObjectAttributes{
		StorageTier:  string(resp.StorageTier),
		ContentType:  resp.ContentType,
		ContentMD5:   resp.ContentMd5,
		ETag:         resp.ETag,
		LastModified: lm,
		BucketName:   bucketName,
		Namespace:    namespace,
	}
}

// NewDomainObjectFromAttrs builds a domain.Object from provider-agnostic attributes.
func NewDomainObjectFromAttrs(attrs ObjectAttributes) domain.Object {
	obj := domain.Object{
		Name:        stringValue(attrs.Name),
		Size:        int64Value(attrs.Size),
		StorageTier: attrs.StorageTier,
		ContentType: stringValue(attrs.ContentType),
		ContentMD5:  stringValue(attrs.ContentMD5),
		ETag:        stringValue(attrs.ETag),
		BucketName:  attrs.BucketName,
		Namespace:   attrs.Namespace,
	}

	if attrs.LastModified != nil {
		obj.LastModified = *attrs.LastModified
	}

	return obj
}
