package instance

import (
	"fmt"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// InstanceOutput defines the structure for the JSON output of an instance.
type InstanceOutput struct {
	Name              string                 `json:"Name"`
	ID                string                 `json:"ID"`
	IP                string                 `json:"IP"`
	ImageID           string                 `json:"ImageID"`
	SubnetID          string                 `json:"SubnetID"`
	Shape             string                 `json:"Shape"`
	State             string                 `json:"State"`
	CreatedAt         time.Time              `json:"CreatedAt"`
	Placement         Placement              `json:"Placement"`
	Resources         Resources              `json:"Resources"`
	ImageName         string                 `json:"ImageName,omitempty"`
	ImageOS           string                 `json:"ImageOS,omitempty"`
	InstanceTags      map[string]interface{} `json:"InstanceTags"`
	Hostname          string                 `json:"Hostname,omitempty"`
	SubnetName        string                 `json:"SubnetName,omitempty"`
	VcnID             string                 `json:"VcnID,omitempty"`
	VcnName           string                 `json:"VcnName,omitempty"`
	PrivateDNSEnabled bool                   `json:"PrivateDNSEnabled,omitempty"`
	RouteTableID      string                 `json:"RouteTableID,omitempty"`
	RouteTableName    string                 `json:"RouteTableName,omitempty"`
	SecurityListIDs   []string               `json:"SecurityListIDs,omitempty"`
	SecurityListNames []string               `json:"SecurityListNames,omitempty"`
	NsgIDs            []string               `json:"NsgIDs,omitempty"`
	NsgNames          []string               `json:"NsgNames,omitempty"`
	IsMtlsRequired    *bool                  `json:"IsMtlsRequired,omitempty"`
}

// Placement represents the location of an instance.
type Placement struct {
	Region             string `json:"Region"`
	AvailabilityDomain string `json:"AvailabilityDomain"`
	FaultDomain        string `json:"FaultDomain"`
}

// Resources represent the compute resources of an instance.
type Resources struct {
	VCPUs    int     `json:"VCPUs"`
	MemoryGB float32 `json:"MemoryGB"`
}

// PrintInstancesInfo displays instances in a formatted table or JSON format.
func PrintInstancesInfo(instances []compute.Instance, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showImageDetails bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		outputInstances := make([]InstanceOutput, len(instances))
		for i, inst := range instances {
			outputInstances[i] = InstanceOutput{
				Name:      inst.DisplayName,
				ID:        inst.OCID,
				IP:        inst.PrimaryIP,
				ImageID:   inst.ImageID,
				SubnetID:  inst.SubnetID,
				Shape:     inst.Shape,
				State:     inst.State,
				CreatedAt: inst.TimeCreated,
				Placement: Placement{
					Region:             inst.Region,
					AvailabilityDomain: inst.AvailabilityDomain,
					FaultDomain:        inst.FaultDomain,
				},
				Resources: Resources{
					VCPUs:    inst.VCPUs,
					MemoryGB: inst.MemoryGB,
				},
				ImageName: inst.ImageName,
				ImageOS:   inst.ImageOS,
				InstanceTags: map[string]interface{}{
					"FreeformTags": inst.FreeformTags,
					"DefinedTags":  inst.DefinedTags,
				},
				Hostname:          inst.Hostname,
				SubnetName:        inst.SubnetName,
				VcnID:             inst.VcnID,
				VcnName:           inst.VcnName,
				PrivateDNSEnabled: inst.PrivateDNSEnabled,
				RouteTableID:      inst.RouteTableID,
				RouteTableName:    inst.RouteTableName,
				SecurityListIDs:   inst.SecurityListIDs,
				SecurityListNames: inst.SecurityListNames,
				NsgIDs:            inst.NsgIDs,
				NsgNames:          inst.NsgNames,
				IsMtlsRequired:    inst.IsMtlsRequired,
			}
		}
		return util.MarshalDataToJSONResponse(p, outputInstances, pagination)
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
			instanceData["Hostname"] = instance.Hostname
			instanceData["Private DNS Enabled"] = fmt.Sprintf("%t", instance.PrivateDNSEnabled)
			instanceData["Route Table Name"] = instance.RouteTableName

			// Add Security Lists
			if len(instance.SecurityListNames) > 0 {
				instanceData["Security Lists"] = fmt.Sprintf("%d attached", len(instance.SecurityListNames))
				for i, slName := range instance.SecurityListNames {
					instanceData[fmt.Sprintf("  Security List %d", i+1)] = slName
				}
			} else {
				instanceData["Security Lists"] = "None"
			}

			// Add NSGs
			if len(instance.NsgNames) > 0 {
				instanceData["NSGs"] = fmt.Sprintf("%d attached", len(instance.NsgNames))
				for i, nsgName := range instance.NsgNames {
					instanceData[fmt.Sprintf("  NSG %d", i+1)] = nsgName
				}
			} else {
				instanceData["NSGs"] = "None"
			}

			// Add mTLS requirement
			if instance.IsMtlsRequired != nil {
				instanceData["mTLS Required"] = fmt.Sprintf("%t", *instance.IsMtlsRequired)
			}

			imageKeys := []string{
				"Image Name", "Operating System", "AD", "FD", "Region", "Subnet Name", "VCN Name", "Hostname", "Private DNS Enabled", "Route Table Name", "Security Lists",
			}
			// Add dynamic security list keys
			for i := range instance.SecurityListNames {
				imageKeys = append(imageKeys, fmt.Sprintf("  Security List %d", i+1))
			}
			imageKeys = append(imageKeys, "NSGs")
			// Add dynamic NSG keys
			for i := range instance.NsgNames {
				imageKeys = append(imageKeys, fmt.Sprintf("  NSG %d", i+1))
			}
			// Add mTLS if present
			if instance.IsMtlsRequired != nil {
				imageKeys = append(imageKeys, "mTLS Required")
			}
			orderedKeys = append(orderedKeys, imageKeys...)
		}

		title := util.FormatColoredTitle(appCtx, instance.DisplayName)
		p.PrintKeyValues(title, instanceData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func PrintInstanceInfo(instance *compute.Instance, appCtx *app.ApplicationContext, useJSON bool, showDetails bool) error {
	p := printer.New(appCtx.Stdout)

	if instance == nil {
		return fmt.Errorf("instance is nil")
	}

	if useJSON {
		out := InstanceOutput{
			Name:      instance.DisplayName,
			ID:        instance.OCID,
			IP:        instance.PrimaryIP,
			ImageID:   instance.ImageID,
			SubnetID:  instance.SubnetID,
			Shape:     instance.Shape,
			State:     instance.State,
			CreatedAt: instance.TimeCreated,
			Placement: Placement{
				Region:             instance.Region,
				AvailabilityDomain: instance.AvailabilityDomain,
				FaultDomain:        instance.FaultDomain,
			},
			Resources: Resources{
				VCPUs:    instance.VCPUs,
				MemoryGB: instance.MemoryGB,
			},
			ImageName: instance.ImageName,
			ImageOS:   instance.ImageOS,
			InstanceTags: map[string]interface{}{
				"FreeformTags": instance.FreeformTags,
				"DefinedTags":  instance.DefinedTags,
			},
			Hostname:          instance.Hostname,
			SubnetName:        instance.SubnetName,
			VcnID:             instance.VcnID,
			VcnName:           instance.VcnName,
			PrivateDNSEnabled: instance.PrivateDNSEnabled,
			RouteTableID:      instance.RouteTableID,
			RouteTableName:    instance.RouteTableName,
			SecurityListIDs:   instance.SecurityListIDs,
			SecurityListNames: instance.SecurityListNames,
			NsgIDs:            instance.NsgIDs,
			NsgNames:          instance.NsgNames,
			IsMtlsRequired:    instance.IsMtlsRequired,
		}
		return p.MarshalToJSON(out)
	}

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

	if showDetails {
		instanceData["Image Name"] = instance.ImageName
		instanceData["Operating System"] = instance.ImageOS
		instanceData["AD"] = instance.AvailabilityDomain
		instanceData["FD"] = instance.FaultDomain
		instanceData["Region"] = instance.Region
		instanceData["Subnet Name"] = instance.SubnetName
		instanceData["VCN Name"] = instance.VcnName
		instanceData["Hostname"] = instance.Hostname
		instanceData["Private DNS Enabled"] = fmt.Sprintf("%t", instance.PrivateDNSEnabled)
		instanceData["Route Table Name"] = instance.RouteTableName

		// Add Security Lists
		if len(instance.SecurityListNames) > 0 {
			instanceData["Security Lists"] = fmt.Sprintf("%d attached", len(instance.SecurityListNames))
			for i, slName := range instance.SecurityListNames {
				instanceData[fmt.Sprintf("  Security List %d", i+1)] = slName
			}
		} else {
			instanceData["Security Lists"] = "None"
		}

		// Add NSGs
		if len(instance.NsgNames) > 0 {
			instanceData["NSGs"] = fmt.Sprintf("%d attached", len(instance.NsgNames))
			for i, nsgName := range instance.NsgNames {
				instanceData[fmt.Sprintf("  NSG %d", i+1)] = nsgName
			}
		} else {
			instanceData["NSGs"] = "None"
		}

		// Add mTLS requirement
		if instance.IsMtlsRequired != nil {
			instanceData["mTLS Required"] = fmt.Sprintf("%t", *instance.IsMtlsRequired)
		}

		imageKeys := []string{
			"Image Name", "Operating System", "AD", "FD", "Region", "Subnet Name", "VCN Name", "Hostname", "Private DNS Enabled", "Route Table Name", "Security Lists",
		}
		// Add dynamic security list keys
		for i := range instance.SecurityListNames {
			imageKeys = append(imageKeys, fmt.Sprintf("  Security List %d", i+1))
		}
		imageKeys = append(imageKeys, "NSGs")
		// Add dynamic NSG keys
		for i := range instance.NsgNames {
			imageKeys = append(imageKeys, fmt.Sprintf("  NSG %d", i+1))
		}
		// Add mTLS if present
		if instance.IsMtlsRequired != nil {
			imageKeys = append(imageKeys, "mTLS Required")
		}
		orderedKeys = append(orderedKeys, imageKeys...)
	}

	title := util.FormatColoredTitle(appCtx, instance.DisplayName)
	p.PrintKeyValues(title, instanceData, orderedKeys)

	return nil
}
