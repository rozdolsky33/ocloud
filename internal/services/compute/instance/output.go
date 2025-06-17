package instance

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

// PrintInstancesTable displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintInstancesTable(instances []Instance, appCtx *app.ApplicationContext, pagination *PaginationInfo, useJSON bool, showImageDetails bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return marshalInstancesToJSON(p, instances, pagination)
	}

	// Handle the case where no instances are found.
	if len(instances) == 0 {
		fmt.Fprintln(appCtx.Stdout, "No instances found.") // Write to the context's writer.
		if pagination != nil && pagination.TotalCount > 0 {
			fmt.Fprintf(appCtx.Stdout, "Page %d is empty. Total records: %d\n", pagination.CurrentPage, pagination.TotalCount)
			if pagination.CurrentPage > 1 {
				fmt.Fprintf(appCtx.Stdout, "Try a lower page number (e.g., --page %d)\n", pagination.CurrentPage-1)
			}
		}
		return nil
	}

	// Print each instance as a separate key-value table with a colored title.
	for _, instance := range instances {
		// Create instance data map
		instanceData := map[string]string{
			"ID":         instance.ID,
			"AD":         instance.Placement.AvailabilityDomain,
			"FD":         instance.Placement.FaultDomain,
			"Region":     instance.Placement.Region,
			"Shape":      instance.Shape,
			"vCPUs":      fmt.Sprintf("%d", instance.Resources.VCPUs),
			"Created":    instance.CreatedAt.String(),
			"Subnet ID":  instance.SubnetID,
			"Name":       instance.Name,
			"Private IP": instance.IP,
			"Memory":     fmt.Sprintf("%d GB", int(instance.Resources.MemoryGB)),
			"State":      string(instance.State),
		}

		// Define ordered keys
		orderedKeys := []string{
			"ID", "AD", "FD", "Region", "Shape", "vCPUs",
			"Created", "Subnet ID", "Name", "Private IP", "Memory", "State",
		}

		// Add image details if available
		if showImageDetails && instance.ImageID != "" {
			// Add image ID
			instanceData["Image ID"] = instance.ImageID

			// Add an image name if available
			if instance.ImageDetails.ImageName != "" {
				instanceData["Image Name"] = instance.ImageDetails.ImageName
			}

			// Add an operating system if available
			if instance.ImageDetails.ImageOS != "" {
				instanceData["Operating System"] = instance.ImageDetails.ImageOS
			}

			// Add image details to ordered keys
			imageKeys := []string{
				"Image ID",
				"Image Name",
				"Operating System",
			}

			// Insert image keys after the "State" key
			newOrderedKeys := make([]string, 0, len(orderedKeys)+len(imageKeys))
			for _, key := range orderedKeys {
				newOrderedKeys = append(newOrderedKeys, key)
				if key == "State" {
					newOrderedKeys = append(newOrderedKeys, imageKeys...)
				}
			}
			orderedKeys = newOrderedKeys

			// Add free-form tags if available
			if len(instance.ImageDetails.ImageFreeformTags) > 0 {
				for k, v := range instance.ImageDetails.ImageFreeformTags {
					tagKey := fmt.Sprintf("Image Tag (Free): %s", k)
					instanceData[tagKey] = v
				}
			}

			// Add defined tags if available
			if len(instance.ImageDetails.ImageDefinedTags) > 0 {
				for namespace, tags := range instance.ImageDetails.ImageDefinedTags {
					for k, v := range tags {
						tagKey := fmt.Sprintf("Image Tag (Defined): %s.%s", namespace, k)
						if v != nil {
							instanceData[tagKey] = fmt.Sprintf("%v", v)
						}
					}
				}
			}
		}

		// Create the colored title using components from the app context.
		coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
		coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
		coloredInstance := text.Colors{text.FgBlue}.Sprint(instance.Name)
		title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredInstance)

		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, instanceData, orderedKeys)
	}

	logPaginationInfo(pagination, appCtx)
	return nil
}

// logPaginationInfo logs pagination information if available.
func logPaginationInfo(pagination *PaginationInfo, appCtx *app.ApplicationContext) {
	// Log pagination information if available
	if pagination != nil {
		// Calculate the total records displayed so far
		totalRecordsDisplayed := pagination.CurrentPage * pagination.Limit
		if totalRecordsDisplayed > pagination.TotalCount {
			totalRecordsDisplayed = pagination.TotalCount
		}

		// Log pagination information at the INFO level
		appCtx.Logger.Info("--- Pagination Information ---",
			"page", pagination.CurrentPage,
			"records", fmt.Sprintf("%d/%d", totalRecordsDisplayed, pagination.TotalCount),
			"limit", pagination.Limit)

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

// marshalInstancesToJSON now accepts a printer and returns an error.
func marshalInstancesToJSON(p *printer.Printer, instances []Instance, pagination *PaginationInfo) error {
	response := JSONResponse{
		Instances:  instances,
		Pagination: pagination,
	}
	// Use the printer's method to marshal. It will write to the correct output.
	return p.MarshalToJSON(response)
}
