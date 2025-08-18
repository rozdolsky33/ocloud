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
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[Instance](p, instances, pagination)
	}

	if util.ValidateAndReportEmpty(instances, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each instance as a separate key-value.
	for _, instance := range instances {
		instanceData := map[string]string{
			"Shape":      instance.Shape,
			"vCPUs":      fmt.Sprintf("%d", instance.Resources.VCPUs),
			"Created":    instance.CreatedAt.String(),
			"Name":       instance.Name,
			"Private IP": instance.IP,
			"Memory":     fmt.Sprintf("%d GB", int(instance.Resources.MemoryGB)),
			"State":      string(instance.State),
		}

		orderedKeys := []string{
			"Name", "Shape", "vCPUs", "Memory",
			"Created", "Private IP", "State",
			"Boot Volume State",
		}

		// Add image details if available
		if showImageDetails {

			if instance.ImageOS != "" {
				instanceData["Operating System"] = instance.ImageOS
			}
			if instance.ImageName != "" {
				instanceData["Image Name"] = instance.ImageName
			}

			if instance.Placement.AvailabilityDomain != "" {
				instanceData["AD"] = instance.Placement.AvailabilityDomain
			}

			if instance.Placement.FaultDomain != "" {
				instanceData["FD"] = instance.Placement.FaultDomain
			}
			if instance.Placement.Region != "" {
				instanceData["Region"] = instance.Placement.Region
			}

			if instance.SubnetName != "" {
				instanceData["Subnet Name"] = instance.SubnetName
			}

			if instance.VcnName != "" {
				instanceData["VCN Name"] = instance.VcnName
			}

			if instance.Hostname != "" {
				instanceData["Hostname"] = instance.Hostname
			}

			instanceData["Private DNS Enabled"] = fmt.Sprintf("%t", instance.PrivateDNSEnabled)

			if instance.RouteTableName != "" {
				instanceData["Route Table Name"] = instance.RouteTableName
			}

			// Add image details to ordered keys
			imageKeys := []string{
				"Image Name",
				"Operating System",
				"AD",
				"FD",
				"Region",
				"Subnet Name",
				"VCN Name",
				"Hostname",
				"Private DNS Enabled",
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

		p.PrintKeyValues(title, instanceData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
