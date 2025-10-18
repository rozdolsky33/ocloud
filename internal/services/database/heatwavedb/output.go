package heatwavedb

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintHeatWaveDbInfo prints a single HeatWave DB.
// - useJSON: if true, prints the single DB as JSON (no pagination envelope)
// - showAll: if true, prints the detailed view; otherwise, prints the summary view
func PrintHeatWaveDbInfo(db *database.HeatWaveDatabase, appCtx *app.ApplicationContext, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(db)
	}

	return printOneHeatWaveDb(p, appCtx, db, showAll)
}

// PrintHeatWaveDbsInfo prints a list of HeatWave DBs.
// - pagination: optional, will be adjusted and logged if provided
// - useJSON: if true, prints databases with util.MarshalDataToJSONResponse
// - showAll: if true, prints detailed view; otherwise summary view
func PrintHeatWaveDbsInfo(databases []database.HeatWaveDatabase, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		if len(databases) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
		return util.MarshalDataToJSONResponse[database.HeatWaveDatabase](p, databases, pagination)
	}

	if util.ValidateAndReportEmpty(databases, pagination, appCtx.Stdout) {
		return nil
	}

	for _, db := range databases {
		if err := printOneHeatWaveDb(p, appCtx, &db, showAll); err != nil {
			return err
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func printOneHeatWaveDb(p *printer.Printer, appCtx *app.ApplicationContext, db *database.HeatWaveDatabase, showAll bool) error {
	title := util.FormatColoredTitle(appCtx, db.DisplayName)

	// Prefer names to IDs when available
	subnetVal := db.SubnetId
	if db.SubnetName != "" {
		subnetVal = db.SubnetName
	}
	vcnVal := db.VcnID
	if db.VcnName != "" {
		vcnVal = db.VcnName
	}
	nsgVal := ""
	if len(db.NsgNames) > 0 {
		nsgVal = fmt.Sprintf("%v", db.NsgNames)
	} else if len(db.NsgIds) > 0 {
		nsgVal = fmt.Sprintf("%v", db.NsgIds)
	}

	// Storage formatting
	storage := ""
	if db.DataStorageSizeInGBs != nil {
		storage = fmt.Sprintf("%d GB", *db.DataStorageSizeInGBs)
	}

	// HeatWave cluster info
	heatwaveCluster := "No"
	if db.IsHeatWaveClusterAttached != nil && *db.IsHeatWaveClusterAttached {
		heatwaveCluster = "Yes"
		if db.HeatWaveCluster != nil && db.HeatWaveCluster.ClusterSize != nil {
			heatwaveCluster = fmt.Sprintf("Yes (%d nodes)", *db.HeatWaveCluster.ClusterSize)
		}
	}

	// High availability
	highAvailability := boolToString(db.IsHighlyAvailable)

	if !showAll {
		// Summary view
		summary := map[string]string{
			"Lifecycle State":   db.LifecycleState,
			"MySQL Version":     db.MysqlVersion,
			"Shape":             db.ShapeName,
			"Storage":           storage,
			"High Availability": highAvailability,
			"HeatWave Cluster":  heatwaveCluster,
			"Database Mode":     db.DatabaseMode,
			"Access Mode":       db.AccessMode,
			"Subnet":            subnetVal,
			"VCN":               vcnVal,
			"NSGs":              nsgVal,
		}
		if db.TimeCreated != nil {
			summary["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
		}
		ordered := []string{
			"Lifecycle State", "MySQL Version", "Shape", "Storage", "High Availability",
			"HeatWave Cluster", "Database Mode", "Access Mode", "Subnet", "VCN", "NSGs", "Time Created",
		}
		p.PrintKeyValues(title, summary, ordered)
		return nil
	}

	// Detailed view
	details := make(map[string]string)
	orderedKeys := []string{}

	// General
	details["Lifecycle State"] = db.LifecycleState
	details["MySQL Version"] = db.MysqlVersion
	details["Description"] = db.Description
	if db.TimeCreated != nil {
		details["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
	}
	if db.TimeUpdated != nil {
		details["Time Updated"] = db.TimeUpdated.Format("2006-01-02 15:04:05")
	}
	orderedKeys = append(orderedKeys, "Lifecycle State", "MySQL Version", "Description", "Time Created", "Time Updated")

	// Capacity
	details["Shape"] = db.ShapeName
	details["Storage"] = storage
	details["High Availability"] = highAvailability
	orderedKeys = append(orderedKeys, "Shape", "Storage", "High Availability")

	// HeatWave Cluster
	details["HeatWave Cluster"] = heatwaveCluster
	orderedKeys = append(orderedKeys, "HeatWave Cluster")
	if db.HeatWaveCluster != nil {
		if db.HeatWaveCluster.ShapeName != nil {
			details["HeatWave Shape"] = *db.HeatWaveCluster.ShapeName
			orderedKeys = append(orderedKeys, "HeatWave Shape")
		}
		if db.HeatWaveCluster.LifecycleState != "" {
			details["HeatWave State"] = string(db.HeatWaveCluster.LifecycleState)
			orderedKeys = append(orderedKeys, "HeatWave State")
		}
	}

	// Configuration & Mode
	details["Database Mode"] = db.DatabaseMode
	details["Access Mode"] = db.AccessMode
	details["Configuration ID"] = db.ConfigurationId
	orderedKeys = append(orderedKeys, "Database Mode", "Access Mode", "Configuration ID")

	// Network
	details["Subnet"] = subnetVal
	details["VCN"] = vcnVal
	details["NSGs"] = nsgVal
	orderedKeys = append(orderedKeys, "Subnet", "VCN", "NSGs")

	// Endpoints
	if len(db.Endpoints) > 0 {
		for i, endpoint := range db.Endpoints {
			if endpoint.IpAddress != nil {
				key := fmt.Sprintf("Endpoint %d IP", i+1)
				details[key] = *endpoint.IpAddress
				orderedKeys = append(orderedKeys, key)
			}
			if endpoint.Hostname != nil {
				key := fmt.Sprintf("Endpoint %d Hostname", i+1)
				details[key] = *endpoint.Hostname
				orderedKeys = append(orderedKeys, key)
			}
		}
	}

	// Placement
	details["Availability Domain"] = db.AvailabilityDomain
	details["Fault Domain"] = db.FaultDomain
	orderedKeys = append(orderedKeys, "Availability Domain", "Fault Domain")

	// Backup Policy
	if db.BackupPolicy != nil {
		if db.BackupPolicy.IsEnabled != nil {
			details["Backup Enabled"] = boolToString(db.BackupPolicy.IsEnabled)
			orderedKeys = append(orderedKeys, "Backup Enabled")
		}
		if db.BackupPolicy.RetentionInDays != nil {
			details["Backup Retention"] = fmt.Sprintf("%d days", *db.BackupPolicy.RetentionInDays)
			orderedKeys = append(orderedKeys, "Backup Retention")
		}
	}

	p.PrintKeyValues(title, details, orderedKeys)
	return nil
}

//-------------------------------------------------Helpers--------------------------------------------------------------

func boolToString(v *bool) string {
	if v == nil {
		return ""
	}
	if *v {
		return "true"
	}
	return "false"
}
