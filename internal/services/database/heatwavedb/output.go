package heatwavedb

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintHeatWaveDbInfo prints a single HeatWave DB.
func PrintHeatWaveDbInfo(db *database.HeatWaveDatabase, appCtx *app.ApplicationContext, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(db)
	}

	return printOneHeatWaveDb(p, appCtx, db, showAll)
}

// PrintHeatWaveDbsInfo prints a list of HeatWave DBs.
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
	} else {
		nsgVal = "No"
	}

	// Storage formatting - prefer DataStorage object if available
	var storage, allocatedStorage, storageLimit string
	var autoExpandEnabled bool

	if db.DataStorage != nil {
		if db.DataStorage.DataStorageSizeInGBs != nil {
			sizeInBytes := int64(*db.DataStorage.DataStorageSizeInGBs) * 1024 * 1024 * 1024
			storage = util.HumanizeBytesIEC(sizeInBytes)
		}
		if db.DataStorage.AllocatedStorageSizeInGBs != nil {
			allocatedInBytes := int64(*db.DataStorage.AllocatedStorageSizeInGBs) * 1024 * 1024 * 1024
			allocatedStorage = util.HumanizeBytesIEC(allocatedInBytes)
		}
		if db.DataStorage.DataStorageSizeLimitInGBs != nil {
			limitInBytes := int64(*db.DataStorage.DataStorageSizeLimitInGBs) * 1024 * 1024 * 1024
			storageLimit = util.HumanizeBytesIEC(limitInBytes)
		}
		if db.DataStorage.IsAutoExpandStorageEnabled != nil {
			autoExpandEnabled = *db.DataStorage.IsAutoExpandStorageEnabled
		}
	} else if db.DataStorageSizeInGBs != nil {
		sizeInBytes := int64(*db.DataStorageSizeInGBs) * 1024 * 1024 * 1024
		storage = util.HumanizeBytesIEC(sizeInBytes)
	}

	heatwaveCluster := "No"
	if db.IsHeatWaveClusterAttached != nil && *db.IsHeatWaveClusterAttached {
		heatwaveCluster = "Yes"
		if db.HeatWaveCluster != nil && db.HeatWaveCluster.ClusterSize != nil {
			heatwaveCluster = fmt.Sprintf("Yes (%d nodes)", *db.HeatWaveCluster.ClusterSize)
		}
	}

	highAvailability := boolToString(db.IsHighlyAvailable)

	primaryIP := db.IpAddress
	primaryPort := ""
	if db.Port != nil {
		primaryPort = fmt.Sprintf("%d", *db.Port)
	}
	primaryFQDN := ""
	if db.HostnameLabel != "" && db.SubnetName != "" && db.VcnName != "" {
		// Construct FQDN pattern: hostname.subnet.vcn.oraclevcn.com
		primaryFQDN = fmt.Sprintf("%s.%s.%s.oraclevcn.com", db.HostnameLabel, db.SubnetName, db.VcnName)
	}

	if !showAll {
		// Summary view - Essential operational info
		summary := map[string]string{
			"Lifecycle State":   db.LifecycleState,
			"MySQL Version":     db.MysqlVersion,
			"Shape":             db.ShapeName,
			"High Availability": highAvailability,
			"HeatWave Cluster":  heatwaveCluster,
			"Database Mode":     db.DatabaseMode,
			"Access Mode":       db.AccessMode,
			"Private IP":        primaryIP,
			"Port":              primaryPort,
			"Subnet":            subnetVal,
			"VCN":               vcnVal,
		}

		// Add ECPU and memory info if shape details are available
		if ecpu, memory, found := getMySQLShapeDetails(db.ShapeName); found {
			summary["ECPUs"] = fmt.Sprintf("%d", ecpu)
			summary["Memory"] = fmt.Sprintf("%d GB", memory)
		}

		// Add storage details
		if storage != "" {
			summary["Storage Used"] = storage
		}
		if allocatedStorage != "" && allocatedStorage != storage {
			summary["Storage Allocated"] = allocatedStorage
		}
		if storageLimit != "" {
			summary["Storage Limit"] = storageLimit
		}
		if db.DataStorage != nil && db.DataStorage.IsAutoExpandStorageEnabled != nil {
			summary["Auto-Expand Storage"] = boolToString(db.DataStorage.IsAutoExpandStorageEnabled)
		}

		if db.TimeCreated != nil {
			summary["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
		}

		ordered := []string{
			"Lifecycle State", "MySQL Version", "Shape", "ECPUs", "Memory",
			"Storage Used", "Storage Allocated", "Storage Limit", "Auto-Expand Storage",
			"High Availability", "HeatWave Cluster", "Database Mode", "Access Mode",
			"Private IP", "Port", "Subnet", "VCN", "Time Created",
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
	// Only show description if not empty
	if db.Description != "" {
		details["Description"] = db.Description
	}
	if db.TimeCreated != nil {
		details["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
	}
	if db.TimeUpdated != nil {
		details["Time Updated"] = db.TimeUpdated.Format("2006-01-02 15:04:05")
	}

	generalKeys := []string{"Lifecycle State", "MySQL Version"}
	if db.Description != "" {
		generalKeys = append(generalKeys, "Description")
	}
	generalKeys = append(generalKeys, "Time Created", "Time Updated")
	orderedKeys = append(orderedKeys, generalKeys...)

	// Capacity
	details["Shape"] = db.ShapeName
	if ecpu, memory, found := getMySQLShapeDetails(db.ShapeName); found {
		details["ECPUs"] = fmt.Sprintf("%d", ecpu)
		details["Memory"] = fmt.Sprintf("%d GB", memory)
	}
	details["High Availability"] = highAvailability
	orderedKeys = append(orderedKeys, "Shape", "ECPUs", "Memory", "High Availability")

	// Storage - consolidate when used == allocated
	if storage != "" && allocatedStorage != "" && storage != allocatedStorage {
		// Show both when different (expansion occurred)
		details["Storage Used"] = storage
		details["Storage Allocated"] = allocatedStorage
		orderedKeys = append(orderedKeys, "Storage Used", "Storage Allocated")
	} else if storage != "" {
		details["Storage"] = storage
		orderedKeys = append(orderedKeys, "Storage")
	}

	// Storage limit (show max capacity)
	if storageLimit != "" {
		details["Storage Limit"] = storageLimit
		orderedKeys = append(orderedKeys, "Storage Limit")
	}

	// Auto-expand settings (only show if enabled or if max size is set)
	if db.DataStorage != nil {
		if autoExpandEnabled {
			details["Auto-Expand"] = "Enabled"
			orderedKeys = append(orderedKeys, "Auto-Expand")
			if db.DataStorage.MaxStorageSizeInGBs != nil {
				maxInBytes := int64(*db.DataStorage.MaxStorageSizeInGBs) * 1024 * 1024 * 1024
				details["Max Expand Size"] = util.HumanizeBytesIEC(maxInBytes)
				orderedKeys = append(orderedKeys, "Max Expand Size")
			}
		}
	}

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
	details["Crash Recovery"] = db.CrashRecovery
	orderedKeys = append(orderedKeys, "Database Mode", "Access Mode", "Crash Recovery")

	details["Subnet"] = subnetVal
	details["VCN"] = vcnVal
	if nsgVal != "" && nsgVal != "No" {
		details["NSGs"] = nsgVal
		orderedKeys = append(orderedKeys, "Subnet", "VCN", "NSGs")
	} else {
		orderedKeys = append(orderedKeys, "Subnet", "VCN")
	}

	// Connection Endpoint - Consolidate connection info
	connectionInfo := primaryIP
	if primaryPort != "" {
		connectionInfo = fmt.Sprintf("%s:%s", primaryIP, primaryPort)
		if db.PortX != nil {
			connectionInfo = fmt.Sprintf("%s:%s (X: %d)", primaryIP, primaryPort, *db.PortX)
		}
	}
	details["Endpoint"] = connectionInfo
	if primaryFQDN != "" {
		details["Internal FQDN"] = primaryFQDN
		orderedKeys = append(orderedKeys, "Endpoint", "Internal FQDN")
	} else {
		orderedKeys = append(orderedKeys, "Endpoint")
	}

	details["Availability Domain"] = db.AvailabilityDomain
	details["Fault Domain"] = db.FaultDomain
	orderedKeys = append(orderedKeys, "Availability Domain", "Fault Domain")

	if db.BackupPolicy != nil {
		if db.BackupPolicy.IsEnabled != nil {
			details["Automatic Backups"] = boolToString(db.BackupPolicy.IsEnabled)
			orderedKeys = append(orderedKeys, "Automatic Backups")
		}
		if db.BackupPolicy.WindowStartTime != nil {
			details["Backup Window"] = *db.BackupPolicy.WindowStartTime
			orderedKeys = append(orderedKeys, "Backup Window")
		}
		if db.BackupPolicy.RetentionInDays != nil {
			details["Retention Days"] = fmt.Sprintf("%d", *db.BackupPolicy.RetentionInDays)
			orderedKeys = append(orderedKeys, "Retention Days")
		}
	}

	if db.PointInTimeRecoveryDetails != nil && db.PointInTimeRecoveryDetails.TimeEarliestRecoveryPoint != nil {
		pitrKeys := []string{}
		details["PITR"] = "Enabled"
		pitrKeys = append(pitrKeys, "PITR")

		if db.PointInTimeRecoveryDetails.TimeEarliestRecoveryPoint != nil {
			details["Earliest Recovery"] = db.PointInTimeRecoveryDetails.TimeEarliestRecoveryPoint.Format("2006-01-02 15:04:05")
			pitrKeys = append(pitrKeys, "Earliest Recovery")
		}
		if db.PointInTimeRecoveryDetails.TimeLatestRecoveryPoint != nil {
			details["Latest Recovery"] = db.PointInTimeRecoveryDetails.TimeLatestRecoveryPoint.Format("2006-01-02 15:04:05")
			pitrKeys = append(pitrKeys, "Latest Recovery")
		}
		orderedKeys = append(orderedKeys, pitrKeys...)
	}

	if db.DeletionPolicy != nil {
		if db.DeletionPolicy.IsDeleteProtected != nil {
			details["Delete Protected"] = boolToString(db.DeletionPolicy.IsDeleteProtected)
			orderedKeys = append(orderedKeys, "Delete Protected")
		}
		if db.DeletionPolicy.FinalBackup != "" {
			details["Final Backup"] = string(db.DeletionPolicy.FinalBackup)
			orderedKeys = append(orderedKeys, "Final Backup")
		}
	}

	if db.MaintenanceInfo != nil && db.MaintenanceInfo.WindowStartTime != nil {
		details["Maintenance Window"] = *db.MaintenanceInfo.WindowStartTime
		orderedKeys = append(orderedKeys, "Maintenance Window")
	}

	if db.EncryptData != nil {
		if db.EncryptData.KeyId != nil {
			details["Encryption Key"] = *db.EncryptData.KeyId
		} else {
			details["Encryption Key"] = "Oracle-managed key"
		}
		orderedKeys = append(orderedKeys, "Encryption Key")
	}

	if db.SecureConnections != nil {
		if db.SecureConnections.CertificateGenerationType != "" {
			details["Security Certificate"] = string(db.SecureConnections.CertificateGenerationType)
			orderedKeys = append(orderedKeys, "Security Certificate")
		}
	}

	if db.ReadEndpoint != nil {
		if db.ReadEndpoint.ReadEndpointIpAddress != nil {
			details["Read Endpoint IP"] = *db.ReadEndpoint.ReadEndpointIpAddress
			orderedKeys = append(orderedKeys, "Read Endpoint IP")
		} else {
			details["Read Endpoint"] = "Disabled"
			orderedKeys = append(orderedKeys, "Read Endpoint")
		}
	}

	if db.DatabaseManagement != "" {
		details["Database Management"] = db.DatabaseManagement
		orderedKeys = append(orderedKeys, "Database Management")
	}

	if len(db.CustomerContacts) > 0 {
		var contacts []string
		for _, contact := range db.CustomerContacts {
			if contact.Email != nil {
				contacts = append(contacts, *contact.Email)
			}
		}
		if len(contacts) > 0 {
			details["Customer Contacts"] = fmt.Sprintf("%v", contacts)
			orderedKeys = append(orderedKeys, "Customer Contacts")
		}
	}

	// useful for troubleshooting
	if db.LifecycleDetails != "" {
		details["Lifecycle Details"] = db.LifecycleDetails
		orderedKeys = append(orderedKeys, "Lifecycle Details")
	}

	p.PrintKeyValues(title, details, orderedKeys)
	return nil
}

//-------------------------------------------------Helpers--------------------------------------------------------------

// getMySQLShapeDetails returns ECPU count and memory in GB for a given MySQL shape.
// MySQL ECPU shapes follow the pattern where memory = ECPU count × 8 GB.
// Returns (ecpuCount, memoryGB, found)
func getMySQLShapeDetails(shapeName string) (int, int, bool) {
	// Map of MySQL shapes to their ECPU and memory specifications
	// Based on Oracle Cloud Infrastructure MySQL HeatWave ECPU shapes (Table 5-1)
	// Source: https://docs.oracle.com/en-us/iaas/mysql-database/doc/supported-shapes.html
	// All shapes follow the 8 GB per ECPU ratio
	shapeSpecs := map[string]struct {
		ecpu   int
		memory int
	}{
		"MySQL.Free": {ecpu: 2, memory: 8},
		"MySQL.2":    {ecpu: 2, memory: 16},
		"MySQL.4":    {ecpu: 4, memory: 32},
		"MySQL.8":    {ecpu: 8, memory: 64},
		"MySQL.16":   {ecpu: 16, memory: 128},
		"MySQL.32":   {ecpu: 32, memory: 256},
		"MySQL.48":   {ecpu: 48, memory: 384},
		"MySQL.64":   {ecpu: 64, memory: 512},
		"MySQL.256":  {ecpu: 256, memory: 1024},
	}

	if spec, exists := shapeSpecs[shapeName]; exists {
		return spec.ecpu, spec.memory, true
	}
	return 0, 0, false
}

func boolToString(v *bool) string {
	if v == nil {
		return ""
	}
	if *v {
		return "true"
	}
	return "false"
}
