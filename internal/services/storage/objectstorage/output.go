package objectstorage

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintBucketsInfo(buckets []objectstorage.Bucket, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
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
			"Name":      bucket.Name,
			"Namespace": bucket.Namespace,
			"Created":   bucket.TimeCreated.String(),
		}

		orderedKeys := []string{"Name", "Namespace", "Created"}

		if showAll {
			bucketData["OCID"] = bucket.OCID
			bucketData["StorageTier"] = bucket.StorageTier
			bucketData["Visibility"] = bucket.Visibility
			bucketData["Encryption"] = bucket.Encryption
			bucketData["Versioning"] = bucket.Versioning
			bucketData["ReplicationEnabled"] = fmt.Sprintf("%v", bucket.ReplicationEnabled)
			bucketData["ReadOnly"] = fmt.Sprintf("%v", bucket.IsReadOnly)
			bucketData["ApproximateCount"] = fmt.Sprintf("%d", bucket.ApproximateCount)
			bucketData["ApproximateSize"] = util.HumanizeBytesIEC(bucket.ApproximateSize)
			bucketData["ApproximateSizeBytes"] = fmt.Sprintf("%d", bucket.ApproximateSize)

			orderedKeys = []string{
				"Name", "OCID", "Namespace", "Created", "StorageTier", "Visibility", "Encryption", "Versioning", "ReplicationEnabled", "ReadOnly", "ApproximateCount", "ApproximateSize", "ApproximateSizeBytes",
			}
		}

		title := util.FormatColoredTitle(appCtx, bucket.Name)
		p.PrintKeyValues(title, bucketData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
