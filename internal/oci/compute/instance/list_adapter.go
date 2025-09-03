package instance

import (
	"fmt"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

// NewImageListModel builds a TUI list for images.
func NewImageListModel(instances []domain.Instance) listx.Model {
	return listx.NewModel("Instances", instances, func(inst domain.Instance) listx.ResourceItemData {
		return listx.ResourceItemData{
			ID:          inst.OCID,
			Title:       inst.DisplayName,
			Description: description(inst),
		}
	})
}

func description(inst domain.Instance) string {
	date := inst.TimeCreated.Format("2006-01-02")

	adfd := inst.AvailabilityDomain
	if inst.FaultDomain != "" {
		adfd = fmt.Sprintf("%s/FD-%s", adfd, inst.FaultDomain)
	}

	cpuMem := ""
	if inst.VCPUs > 0 || inst.MemoryGB > 0 {
		cpuMem = fmt.Sprintf(" %dvCPU/%.1fGB", inst.VCPUs, inst.MemoryGB)
	}

	spec := inst.Shape
	if spec != "" {
		spec += cpuMem
	} else {
		spec = strings.TrimSpace(cpuMem)
	}

	return fmt.Sprintf("%s • %s • %s %s • %s",
		inst.State, spec, inst.Region, adfd, date)
}
