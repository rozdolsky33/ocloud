package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintImagesInfo prints information for a slice of images either as ordered key/value sections or as a JSON response.
// 
// If a pagination object is provided, pagination settings are adjusted and pagination metadata is logged after output.
// When useJSON is true, the images are marshaled and written as a JSON response; otherwise each image is rendered with
// the fields Name, Created, OperatingSystem, OS Version, and LaunchMode. If the result set is empty the function exits
// without producing output.
//
// It returns an error if writing or marshaling the output fails.
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
// When useJSON is true, the image is serialized to JSON to stdout; otherwise the image's fields are rendered as ordered key-value pairs with a colored title.
func PrintImageInfo(image *compute.Image, appCtx *app.ApplicationContext, useJSON bool) error {
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
