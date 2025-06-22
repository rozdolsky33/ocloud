package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintInstancesInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintInstancesInfo(instances []Instance, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showImageDetails bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return util.MarshalDataToJSON[Instance](p, instances, pagination)
	}

	// Handle the case where no instances are found.
	if util.ValidateAndReportEmpty(instances, pagination, appCtx.Stdout) {
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

		title := util.FormatColoredTitle(appCtx, instance.Name)

		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, instanceData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
