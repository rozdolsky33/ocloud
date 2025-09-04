package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintImagesInfo displays instances in a formatted table or JSON format.
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
			"Name":            image.DisplayName,
			"Created":         image.TimeCreated.String(),
			"OS Version":      image.OperatingSystemVersion,
			"OperatingSystem": image.OperatingSystem,
			"LaunchMode":      image.LaunchMode,
		}

		orderedKeys := []string{
			"Name", "Created", "OperatingSystem", "OS Version", "LaunchMode",
		}

		title := util.FormatColoredTitle(appCtx, image.DisplayName)

		p.PrintKeyValues(title, imageData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// PrintImageInfo prints a detailed view of an image.
func PrintImageInfo(image *domain.Image, appCtx *app.ApplicationContext, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(image)
	}

	imageData := map[string]string{
		"ID":              image.OCID,
		"Name":            image.DisplayName,
		"Created":         image.TimeCreated.String(),
		"OS Version":      image.OperatingSystemVersion,
		"OperatingSystem": image.OperatingSystem,
		"LaunchMode":      image.LaunchMode,
	}

	orderedKeys := []string{
		"ID", "Name", "Created", "OperatingSystem", "OS Version", "LaunchMode",
	}

	title := util.FormatColoredTitle(appCtx, image.DisplayName)

	p.PrintKeyValues(title, imageData, orderedKeys)

	return nil
}
