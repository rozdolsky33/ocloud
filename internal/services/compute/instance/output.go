package instance

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintInstancesInfo displays instances in a formatted table or JSON format.
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

	for _, instance := range instances {
		instanceData := map[string]string{
			"Name":       instance.DisplayName,
			"Shape":      instance.Shape,
			"vCPUs":      fmt.Sprintf("%d", instance.VCPUs),
			"Memory":     fmt.Sprintf("%d GB", int(instance.MemoryGB)),
			"Created":    instance.TimeCreated.String(),
			"Private IP": instance.PrimaryIP,
			"State":      instance.State,
		}

		orderedKeys := []string{
			"Name", "Shape", "vCPUs", "Memory", "Created", "Private IP", "State",
		}

		if showImageDetails {
			instanceData["Image Name"] = instance.ImageName
			instanceData["Operating System"] = instance.ImageOS
			instanceData["AD"] = instance.AvailabilityDomain
			instanceData["FD"] = instance.FaultDomain
			instanceData["Region"] = instance.Region
			instanceData["Subnet Name"] = instance.SubnetName
			instanceData["VCN Name"] = instance.VcnName

			imageKeys := []string{
				"Image Name", "Operating System", "AD", "FD", "Region", "Subnet Name", "VCN Name",
			}
			orderedKeys = append(orderedKeys, imageKeys...)
		}

		title := util.FormatColoredTitle(appCtx, instance.DisplayName)
		p.PrintKeyValues(title, instanceData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
