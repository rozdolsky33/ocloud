package autonomousdb

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewDatabaseListModel builds a TUI list for ADBs.
func NewDatabaseListModel(adbs []domain.AutonomousDatabase) tui.Model {
	return tui.NewModel("Autonomous Databases", adbs, func(adb domain.AutonomousDatabase) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          adb.ID,
			Title:       adb.Name,
			Description: describeAutonomousDatabase(adb),
		}
	})
}

func describeAutonomousDatabase(adb domain.AutonomousDatabase) string {
	wv := strings.TrimSpace(strings.Join(
		filterNonEmpty(adb.DbWorkload, adb.DbVersion),
		" ",
	))
	cpu := ""
	if adb.EcpuCount != nil && *adb.EcpuCount > 0 {
		cpu = fmt.Sprintf("%s eCPU", trimFloat(*adb.EcpuCount))
	} else if adb.OcpuCount != nil && *adb.OcpuCount > 0 {
		cpu = fmt.Sprintf("%s OCPU", trimFloat(*adb.OcpuCount))
	} else if adb.CpuCoreCount != nil && *adb.CpuCoreCount > 0 {
		cpu = fmt.Sprintf("%d CPU", *adb.CpuCoreCount)
	}
	storage := ""
	switch {
	case adb.DataStorageSizeInTBs != nil && *adb.DataStorageSizeInTBs > 0:
		storage = fmt.Sprintf("%dTB", *adb.DataStorageSizeInTBs)
	case adb.DataStorageSizeInGBs != nil && *adb.DataStorageSizeInGBs > 0:
		storage = fmt.Sprintf("%dGB", *adb.DataStorageSizeInGBs)
	}

	spec := strings.TrimSpace(strings.Join(filterNonEmpty(cpu, storage), "/"))
	access := ""
	if adb.PrivateEndpointLabel != "" {
		access = "Private " + adb.PrivateEndpointLabel
	} else if adb.SubnetName != "" {
		access = "Private " + adb.SubnetName
	} else {
		access = "Private"
	}

	license := ""
	switch strings.ToUpper(adb.LicenseModel) {
	case "BRING_YOUR_OWN_LICENSE":
		license = "BYOL"
	case "LICENSE_INCLUDED":
		license = "LI"
	}

	auto := ""
	autoFlags := []string{}
	if isTrue(adb.IsAutoScalingEnabled) {
		autoFlags = append(autoFlags, "CPU-auto")
	}
	if isTrue(adb.IsStorageAutoScalingEnabled) {
		autoFlags = append(autoFlags, "Storage-auto")
	}
	if len(autoFlags) > 0 {
		auto = strings.Join(autoFlags, ",")
	}

	dg := ""
	if isTrue(adb.IsDataGuardEnabled) && adb.Role != "" {
		dg = "DG " + strings.ToUpper(adb.Role)
	}

	date := ""
	if adb.TimeCreated != nil && !adb.TimeCreated.IsZero() {
		date = adb.TimeCreated.Format("2006-01-02")
	}

	parts := []string{}
	if adb.LifecycleState != "" {
		parts = append(parts, adb.LifecycleState)
	}
	if wv != "" {
		parts = append(parts, wv)
	}
	if spec != "" {
		parts = append(parts, spec)
	}
	if access != "" {
		parts = append(parts, access)
	}
	if license != "" {
		parts = append(parts, license)
	}
	if dg != "" {
		parts = append(parts, dg)
	}
	if auto != "" {
		parts = append(parts, auto)
	}
	if date != "" {
		parts = append(parts, date)
	}
	return strings.Join(parts, " • ")
}

// --- helpers ---

func isTrue(b *bool) bool { return b != nil && *b }

func filterNonEmpty(vals ...string) []string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}

// trimFloat prints 1 decimal max and trims trailing zeros
func trimFloat(f float32) string {
	s := fmt.Sprintf("%.1f", f)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}
