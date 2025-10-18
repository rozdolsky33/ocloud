package heatwavedb

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewDatabaseListModel builds a TUI list for HeatWave Databases.
func NewDatabaseListModel(dbs []domain.HeatWaveDatabase) tui.Model {
	return tui.NewModel("HeatWave Databases", dbs, func(db domain.HeatWaveDatabase) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          db.ID,
			Title:       db.DisplayName,
			Description: describeHeatWaveDatabase(db),
		}
	})
}

func describeHeatWaveDatabase(db domain.HeatWaveDatabase) string {
	// MySQL version and shape
	versionShape := strings.TrimSpace(strings.Join(
		filterNonEmpty(db.MysqlVersion, db.ShapeName),
		" ",
	))

	// Storage - prefer DataStorage object
	storage := ""
	if db.DataStorage != nil && db.DataStorage.DataStorageSizeInGBs != nil {
		sizeGB := *db.DataStorage.DataStorageSizeInGBs
		if sizeGB >= 1024 {
			storage = fmt.Sprintf("%.1fTB", float64(sizeGB)/1024.0)
		} else {
			storage = fmt.Sprintf("%dGB", sizeGB)
		}
	} else if db.DataStorageSizeInGBs != nil {
		sizeGB := *db.DataStorageSizeInGBs
		if sizeGB >= 1024 {
			storage = fmt.Sprintf("%.1fTB", float64(sizeGB)/1024.0)
		} else {
			storage = fmt.Sprintf("%dGB", sizeGB)
		}
	}

	// High availability
	ha := ""
	if isTrue(db.IsHighlyAvailable) {
		ha = "HA"
	}

	// HeatWave cluster
	heatwave := ""
	if isTrue(db.IsHeatWaveClusterAttached) {
		if db.HeatWaveCluster != nil && db.HeatWaveCluster.ClusterSize != nil {
			heatwave = fmt.Sprintf("HeatWave(%d)", *db.HeatWaveCluster.ClusterSize)
		} else {
			heatwave = "HeatWave"
		}
	}

	// Network - subnet name or VCN name
	network := ""
	if db.SubnetName != "" {
		network = db.SubnetName
	} else if db.VcnName != "" {
		network = db.VcnName
	}

	// Date created
	date := ""
	if db.TimeCreated != nil && !db.TimeCreated.IsZero() {
		date = db.TimeCreated.Format("2006-01-02")
	}

	// Database and access mode
	mode := ""
	if db.DatabaseMode != "" {
		mode = db.DatabaseMode
	}

	// Build description parts
	parts := []string{}
	if db.LifecycleState != "" {
		parts = append(parts, db.LifecycleState)
	}
	if versionShape != "" {
		parts = append(parts, versionShape)
	}
	if storage != "" {
		parts = append(parts, storage)
	}
	if ha != "" {
		parts = append(parts, ha)
	}
	if heatwave != "" {
		parts = append(parts, heatwave)
	}
	if network != "" {
		parts = append(parts, network)
	}
	if mode != "" {
		parts = append(parts, mode)
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
