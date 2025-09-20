package instance

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
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
	cpuMem := ""
	if inst.VCPUs > 0 || inst.MemoryGB > 0 {
		mem := fmt.Sprintf("%.1f", inst.MemoryGB)
		mem = strings.TrimSuffix(mem, ".0")
		cpuMem = fmt.Sprintf(" %dvCPU/%sGB", inst.VCPUs, mem)
	}
	spec := strings.TrimSpace(inst.Shape + cpuMem)
	fd := inst.FaultDomain
	if strings.HasPrefix(fd, "FAULT-DOMAIN-") {
		fd = "FD-" + strings.TrimPrefix(fd, "FAULT-DOMAIN-")
	}
	var locParts []string
	if inst.Region != "" {
		locParts = append(locParts, inst.Region)
	}
	if inst.AvailabilityDomain != "" && fd != "" {
		locParts = append(locParts, inst.AvailabilityDomain+"/"+fd)
	} else if inst.AvailabilityDomain != "" {
		locParts = append(locParts, inst.AvailabilityDomain)
	} else if fd != "" {
		locParts = append(locParts, fd)
	}
	loc := strings.Join(locParts, " ")

	date := ""
	if !inst.TimeCreated.IsZero() {
		date = inst.TimeCreated.Format("2006-01-02")
	}
	parts := make([]string, 0, 4)
	if inst.State != "" {
		parts = append(parts, inst.State)
	}
	if spec != "" {
		parts = append(parts, spec)
	}
	if loc != "" {
		parts = append(parts, loc)
	}
	if date != "" {
		parts = append(parts, date)
	}
	return strings.Join(parts, " • ")
}
