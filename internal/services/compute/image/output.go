package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintImagesInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintImagesInfo(images []Image, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[Image](p, images, pagination)
	}

	if util.ValidateAndReportEmpty(images, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each image as a separate key-value.
	for _, image := range images {
		imageData := map[string]string{
			"Name":            image.Name,
			"Created":         image.CreatedAt.String(),
			"ImageOSVersion":  image.ImageOSVersion,
			"OperatingSystem": image.OperatingSystem,
			"LunchMode":       image.LunchMode,
		}

		orderedKeys := []string{
			"Name", "Created", "ImageName", "ImageOSVersion", "OperatingSystem", "LunchMode",
		}

		title := util.FormatColoredTitle(appCtx, image.Name)

		p.PrintKeyValues(title, imageData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
