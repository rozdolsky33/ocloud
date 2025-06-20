package instance

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

// PrintInstancesInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintInstancesInfo(instances []Instance, appCtx *app.ApplicationContext, pagination *PaginationInfo, useJSON bool, showImageDetails bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		adjustPaginationInfo(pagination)
	}

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
			"ID", "Name", "Shape", "vCPUs", "Memory",
			"Created", "Subnet ID", "Private IP", "State",
			"Boot Volume ID", "Boot Volume State",
		}

		// Add image details if available
		if showImageDetails && instance.ImageID != "" {
			// Add image ID
			instanceData["Image ID"] = instance.ImageID

			// Add an operating system if available
			if instance.ImageOS != "" {
				instanceData["Operating System"] = instance.ImageOS
			}
			if instance.ImageName != "" {
				instanceData["Image Name"] = instance.ImageName
			}

			//Add AD
			if instance.Placement.AvailabilityDomain != "" {
				instanceData["AD"] = instance.Placement.AvailabilityDomain
			}

			// AD FD
			if instance.Placement.FaultDomain != "" {
				instanceData["FD"] = instance.Placement.FaultDomain
			}
			if instance.Placement.Region != "" {
				instanceData["Region"] = instance.Placement.Region
			}

			// Add subnet details
			if instance.SubnetName != "" {
				instanceData["Subnet Name"] = instance.SubnetName
			}
			if instance.VcnID != "" {
				instanceData["VCN ID"] = instance.VcnID
			}
			if instance.VcnName != "" {
				instanceData["VCN Name"] = instance.VcnName
			}

			// Add hostname
			if instance.Hostname != "" {
				instanceData["Hostname"] = instance.Hostname
			}

			// Add private DNS enabled flag
			instanceData["Private DNS Enabled"] = fmt.Sprintf("%t", instance.PrivateDNSEnabled)

			// Add route table details
			if instance.RouteTableID != "" {
				instanceData["Route Table ID"] = instance.RouteTableID
			}
			if instance.RouteTableName != "" {
				instanceData["Route Table Name"] = instance.RouteTableName
			}

			// Add image details to ordered keys
			imageKeys := []string{
				"Image ID",
				"Image Name",
				"Operating System",
				"AD",
				"FD",
				"Region",
				"Subnet Name",
				"VCN ID",
				"VCN Name",
				"Hostname",
				"Private DNS Enabled",
				"Route Table ID",
				"Route Table Name",
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
func marshalInstancesToJSON(p *printer.Printer, instances []Instance, pagination *PaginationInfo) error {
	response := JSONResponse{
		Instances:  instances,
		Pagination: pagination,
	}
	// Use the printer's method to marshal. It will write to the correct output.
	return p.MarshalToJSON(response)
}
