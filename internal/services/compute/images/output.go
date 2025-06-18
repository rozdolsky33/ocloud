package images

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

// PrintImagesInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintImagesInfo(images []Image, appCtx *app.ApplicationContext, pagination *PaginationInfo, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		adjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return marshalInstancesToJSON(p, images, pagination)
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
			"Created":         image.CreatedAt,
			"ImageName":       image.ImageName,
			"ImageOSVersion":  image.ImageOSVersion,
			"OperatingSystem": image.OperatingSystem,
		}

		// Define ordered keys
		orderedKeys := []string{
			"Name", "ID", "Created", "ImageName", "ImageOSVersion", "OperatingSystem",
		}

		// Create the colored title using components from the app context.
		coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
		coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
		coloredInstance := text.Colors{text.FgBlue}.Sprint(image.Name)
		title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredInstance)
		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, imageData, orderedKeys)
	}

	logPaginationInfo(pagination, appCtx)
	return nil
}

// logPaginationInfo logs pagination information if available.
func logPaginationInfo(pagination *PaginationInfo, appCtx *app.ApplicationContext) {
	// Log pagination information if available
	if pagination != nil {
		// Determine if there's a next page
		hasNextPage := pagination.NextPageToken != ""

		// Log pagination information at the INFO level
		appCtx.Logger.Info("--- Pagination Information ---",
			"page", pagination.CurrentPage,
			"records", fmt.Sprintf("%d", pagination.TotalCount),
			"limit", pagination.Limit,
			"nextPage", hasNextPage)

		// Add debug logs for navigation hints
		if pagination.CurrentPage > 1 {
			logger.LogWithLevel(appCtx.Logger, 2, "Pagination navigation",
				"action", "previous page",
				"page", pagination.CurrentPage-1,
				"limit", pagination.Limit)
		}

		// Check if there are more pages after the current page
		// The most direct way to determine if there are more pages is to check if there's a next page token
		if pagination.NextPageToken != "" {
			logger.LogWithLevel(appCtx.Logger, 2, "Pagination navigation",
				"action", "next page",
				"page", pagination.CurrentPage+1,
				"limit", pagination.Limit)
		}
	}
}

// adjustPaginationInfo adjusts the pagination information to ensure that the total count
// is correctly displayed. It calculates the total records displayed so far and updates
// the TotalCount field of the pagination object to match this value.
func adjustPaginationInfo(pagination *PaginationInfo) {
	// Calculate the total records displayed so far
	totalRecordsDisplayed := pagination.CurrentPage * pagination.Limit
	if totalRecordsDisplayed > pagination.TotalCount {
		totalRecordsDisplayed = pagination.TotalCount
	}

	// Update the total count to match the total records displayed so far
	// This ensures that on page 1 we show 20, on page 2 we show 40, on page 3 we show 60, etc.
	pagination.TotalCount = totalRecordsDisplayed
}

// marshalInstancesToJSON now accepts a printer and returns an error.
func marshalInstancesToJSON(p *printer.Printer, images []Image, pagination *PaginationInfo) error {
	response := JSONResponse{
		Images:     images,
		Pagination: pagination,
	}
	// Use the printer's method to marshal. It will write to the correct output.
	return p.MarshalToJSON(response)
}
