package image

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
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
		return util.MarshalDataToJSON[Image](p, images, pagination)
	}

	// Handle the case where no instances are found.
	if len(images) == 0 {
		fmt.Fprintln(appCtx.Stdout, "No instances found.") // Write to the context's writer.
		if pagination != nil && pagination.TotalCount > 0 {
			fmt.Fprintf(appCtx.Stdout, "Page %d is empty. Total records: %d\n", pagination.CurrentPage, pagination.TotalCount)
			if pagination.CurrentPage > 1 {
				fmt.Fprintf(appCtx.Stdout, "Try a lower page number (e.g., --page %d)\n", pagination.CurrentPage-1)
			}
		}
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
		coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
		coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
		coloredInstance := text.Colors{text.FgBlue}.Sprint(image.Name)
		title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredInstance)
		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, imageData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
