package objectstorage

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintBucketsInfo displays buckets in a formatted table or JSON format.
// If pagination info is provided, it adjusts and logs it.
func PrintBucketsInfo(buckets []Bucket, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[objectstorage.Bucket](p, buckets, pagination)
	}

	if util.ValidateAndReportEmpty(buckets, pagination, appCtx.Stdout) {
		return nil
	}

	for _, bucket := range buckets {
		bucketData := map[string]string{
			"Name":                 bucket.Name,
			"Namespace":            bucket.Namespace,
			"Created":              bucket.TimeCreated.String(),
			"OCID":                 bucket.OCID,
			"StorageTier":          bucket.StorageTier,
			"Visibility":           bucket.Visibility,
			"Encryption":           bucket.Encryption,
			"Versioning":           bucket.Versioning,
			"ReplicationEnabled":   fmt.Sprintf("%v", bucket.ReplicationEnabled),
			"ReadOnly":             fmt.Sprintf("%v", bucket.IsReadOnly),
			"ApproximateCount":     fmt.Sprintf("%d", bucket.ApproximateCount),
			"ApproximateSize":      util.HumanizeBytesIEC(bucket.ApproximateSize),
			"ApproximateSizeBytes": fmt.Sprintf("%d", bucket.ApproximateSize),
		}

		orderedKeys := []string{
			"Name", "OCID", "Namespace", "Created", "StorageTier", "Visibility", "Encryption", "Versioning", "ReplicationEnabled", "ReadOnly", "ApproximateCount", "ApproximateSize", "ApproximateSizeBytes",
		}

		title := util.FormatColoredTitle(appCtx, bucket.Name)
		p.PrintKeyValues(title, bucketData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func PrintBucketInfo(bucket *Bucket, appCtx *app.ApplicationContext, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(bucket)
	}

	bucketData := map[string]string{
		"Name":                 bucket.Name,
		"Namespace":            bucket.Namespace,
		"Created":              bucket.TimeCreated.String(),
		"OCID":                 bucket.OCID,
		"StorageTier":          bucket.StorageTier,
		"Visibility":           bucket.Visibility,
		"Encryption":           bucket.Encryption,
		"Versioning":           bucket.Versioning,
		"ReplicationEnabled":   fmt.Sprintf("%v", bucket.ReplicationEnabled),
		"ReadOnly":             fmt.Sprintf("%v", bucket.IsReadOnly),
		"ApproximateCount":     fmt.Sprintf("%d", bucket.ApproximateCount),
		"ApproximateSize":      util.HumanizeBytesIEC(bucket.ApproximateSize),
		"ApproximateSizeBytes": fmt.Sprintf("%d", bucket.ApproximateSize),
	}

	orderedKeys := []string{
		"Name", "OCID", "Namespace", "Created", "StorageTier", "Visibility", "Encryption", "Versioning", "ReplicationEnabled", "ReadOnly", "ApproximateCount", "ApproximateSize", "ApproximateSizeBytes",
	}

	title := util.FormatColoredTitle(appCtx, bucket.Name)
	p.PrintKeyValues(title, bucketData, orderedKeys)

	return nil
}

// PrintObjectInfo displays object details in a formatted table or JSON format.
func PrintObjectInfo(obj *Object, appCtx *app.ApplicationContext, region string, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(obj)
	}

	// Generate URLs
	encodedName := urlEncode(obj.Name)
	legacyURL := fmt.Sprintf("https://objectstorage.%s.oraclecloud.com/n/%s/b/%s/o/%s",
		region, obj.Namespace, obj.BucketName, encodedName)
	newURL := fmt.Sprintf("https://%s.objectstorage.%s.oci.customer-oci.com/n/%s/b/%s/o/%s",
		obj.Namespace, region, obj.Namespace, obj.BucketName, encodedName)

	objectData := map[string]string{
		"Name":         obj.Name,
		"URL (legacy)": legacyURL,
		"URL (new)":    newURL,
		"StorageTier":  obj.StorageTier,
		"Size":         util.HumanizeBytesIEC(obj.Size),
		"SizeBytes":    fmt.Sprintf("%d", obj.Size),
		"ContentType":  obj.ContentType,
		"ContentMD5":   obj.ContentMD5,
		"ETag":         obj.ETag,
		"LastModified": obj.LastModified.String(),
	}

	orderedKeys := []string{
		"Name", "URL (legacy)", "URL (new)", "StorageTier", "Size", "SizeBytes", "ContentType", "ContentMD5", "ETag", "LastModified",
	}

	title := util.FormatColoredTitle(appCtx, obj.Name)
	p.PrintKeyValues(title, objectData, orderedKeys)

	return nil
}

// urlEncode encodes a string for use in a URL path.
func urlEncode(s string) string {
	// Simple URL encoding for object names
	var result string
	for _, c := range s {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' || c == '/':
			result += string(c)
		default:
			result += fmt.Sprintf("%%%02X", c)
		}
	}
	return result
}
