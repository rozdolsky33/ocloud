package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintImagesInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintImagesInfo(images []Image, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return util.MarshalDataToJSONResponse[Image](p, images, pagination)
	}

	if util.ValidateAndReportEmpty(images, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each image as a separate key-value table with a colored title.
	for _, image := range images {
		// Create image data map
		imageData := map[string]string{
			"Name":            image.Name,
			"ID":              image.ID,
			"Created":         image.CreatedAt.String(),
			"ImageOSVersion":  image.ImageOSVersion,
			"OperatingSystem": image.OperatingSystem,
			"LunchMode":       image.LunchMode,
		}

		// Define ordered keys
		orderedKeys := []string{
			"Name", "ID", "Created", "ImageName", "ImageOSVersion", "OperatingSystem", "LunchMode",
		}

		// Create the colored title using components from the app context.
		title := util.FormatColoredTitle(appCtx, image.Name)

		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, imageData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
