package instance

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

// PrintInstancesTable displays instances in a formatted table or JSON format.
// If useJSON is true, it outputs the instances as JSON, otherwise as a table.
func PrintInstancesTable(instances []Instance, appCtx *app.AppContext, pagination *PaginationInfo, useJSON bool) {
	// If JSON output is requested, print instances as JSON
	if useJSON {
		marshalInstancesToJSON(instances, appCtx, pagination)
		return
	}

	// Create a table printer with the tenancy name as the title
	tablePrinter := printer.NewTablePrinter(appCtx.TenancyName)

	// Convert instances to a format suitable for the printer
	if len(instances) == 0 {
		fmt.Println("No instances found.")
		if pagination != nil && pagination.TotalCount > 0 {
			fmt.Printf("Page %d is empty. Total records: %d\n", pagination.CurrentPage, pagination.TotalCount)
			if pagination.CurrentPage > 1 {
				fmt.Printf("Try a lower page number (e.g., --page %d)\n", pagination.CurrentPage-1)
			}
		}
		return
	}

	// Print each instance as a key-value table with a title
	for _, instance := range instances {
		// Create a map with the instance data
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

		// Define the order of keys to match the example
		orderedKeys := []string{
			"ID",
			"AD",
			"FD",
			"Region",
			"Shape",
			"vCPUs",
			"Created",
			"Subnet ID",
			"Name",
			"Private IP",
			"Memory",
			"State",
		}

		// Print the table with ordered keys and colored title components
		tablePrinter.PrintKeyValueTableWithTitleOrdered(appCtx, instance.Name, instanceData, orderedKeys)
	}

	logPaginationInfo(pagination, appCtx)
}

// logPaginationInfo logs pagination information if available.
func logPaginationInfo(pagination *PaginationInfo, appCtx *app.AppContext) {
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
		if pagination.CurrentPage*pagination.Limit < pagination.TotalCount {
			logger.LogWithLevel(appCtx.Logger, 2, "Pagination navigation",
				"action", "next page",
				"page", pagination.CurrentPage+1,
				"limit", pagination.Limit)
		}
	}
}

// marshalInstancesToJSON marshals instances to JSON and prints the result.
func marshalInstancesToJSON(instances []Instance, appCtx *app.AppContext, pagination *PaginationInfo) {
	response := JSONResponse{
		Instances:  instances,
		Pagination: pagination,
	}

	// Use the printer package to marshal the response to JSON
	printer.MarshalToJSON(response, appCtx)
}
