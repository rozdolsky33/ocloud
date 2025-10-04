package objectstorage

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintBucketsInfo(buckets []objectstorage.Bucket, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
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

		orderedKeys := []string{
			"Name", "Namespace", "Created",
		}

		title := util.FormatColoredTitle(appCtx, bucket.Name)

		p.PrintKeyValues(title, bucketData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
