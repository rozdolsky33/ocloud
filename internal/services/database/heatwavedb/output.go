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
	} else {
		nsgVal = "No"
	}

	// Storage formatting - prefer DataStorage object if available
	var storage, allocatedStorage, storageLimit string
	var autoExpandEnabled bool

	if db.DataStorage != nil {
		// Use the new DataStorage object for detailed information
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
		// Fallback to the deprecated field
		sizeInBytes := int64(*db.DataStorageSizeInGBs) * 1024 * 1024 * 1024
		storage = util.HumanizeBytesIEC(sizeInBytes)
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

	// Primary endpoint info
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
	// Add ECPU and memory info if shape details are available
	if ecpu, memory, found := getMySQLShapeDetails(db.ShapeName); found {
		details["ECPUs"] = fmt.Sprintf("%d", ecpu)
		details["Memory"] = fmt.Sprintf("%d GB", memory)
	}
	details["High Availability"] = highAvailability
	orderedKeys = append(orderedKeys, "Shape", "ECPUs", "Memory", "High Availability")

	// Storage - detailed information
	if storage != "" {
		details["Storage Used"] = storage
		orderedKeys = append(orderedKeys, "Storage Used")
	}
	if allocatedStorage != "" {
		details["Storage Allocated"] = allocatedStorage
		orderedKeys = append(orderedKeys, "Storage Allocated")
	}
	if storageLimit != "" {
		details["Storage Limit"] = storageLimit
		orderedKeys = append(orderedKeys, "Storage Limit")
	}
	if db.DataStorage != nil {
		if db.DataStorage.IsAutoExpandStorageEnabled != nil {
			details["Auto-Expand Storage"] = boolToString(db.DataStorage.IsAutoExpandStorageEnabled)
			orderedKeys = append(orderedKeys, "Auto-Expand Storage")
		}
		if db.DataStorage.MaxStorageSizeInGBs != nil && autoExpandEnabled {
			maxInBytes := int64(*db.DataStorage.MaxStorageSizeInGBs) * 1024 * 1024 * 1024
			details["Max Expand Size"] = util.HumanizeBytesIEC(maxInBytes)
			orderedKeys = append(orderedKeys, "Max Expand Size")
		}
	}

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
	details["Crash Recovery"] = db.CrashRecovery
	orderedKeys = append(orderedKeys, "Database Mode", "Access Mode", "Configuration ID", "Crash Recovery")

	// Network
	details["Subnet"] = subnetVal
	details["Subnet Type"] = "Regional" // HeatWave DB systems use regional subnets
	details["VCN"] = vcnVal
	details["NSGs"] = nsgVal
	orderedKeys = append(orderedKeys, "Subnet", "Subnet Type", "VCN", "NSGs")

	// Primary Endpoint - Critical connection info
	details["Private IP"] = primaryIP
	if primaryPort != "" {
		details["Database Port"] = primaryPort
	}
	if db.PortX != nil {
		details["X Protocol Port"] = fmt.Sprintf("%d", *db.PortX)
	}
	if primaryFQDN != "" {
		details["Internal FQDN"] = primaryFQDN
	}
	orderedKeys = append(orderedKeys, "Private IP", "Database Port", "X Protocol Port", "Internal FQDN")

	// Additional Endpoints
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
			if endpoint.Port != nil {
				key := fmt.Sprintf("Endpoint %d Port", i+1)
				details[key] = fmt.Sprintf("%d", *endpoint.Port)
				orderedKeys = append(orderedKeys, key)
			}
		}
	}

	// Placement
	details["Availability Domain"] = db.AvailabilityDomain
	details["Fault Domain"] = db.FaultDomain
	orderedKeys = append(orderedKeys, "Availability Domain", "Fault Domain")

	// Backup Policy - Critical for SREs
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

	// Point-in-time recovery
	if db.PointInTimeRecoveryDetails != nil {
		details["Point-in-time Recovery"] = "Enabled"
		if db.PointInTimeRecoveryDetails.TimeEarliestRecoveryPoint != nil {
			details["Earliest Recovery Point"] = db.PointInTimeRecoveryDetails.TimeEarliestRecoveryPoint.Format("2006-01-02 15:04:05")
		}
		orderedKeys = append(orderedKeys, "Point-in-time Recovery", "Earliest Recovery Point")
	}

	// Deletion protection
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

	// Maintenance window
	if db.MaintenanceInfo != nil && db.MaintenanceInfo.WindowStartTime != nil {
		details["Maintenance Window"] = *db.MaintenanceInfo.WindowStartTime
		orderedKeys = append(orderedKeys, "Maintenance Window")
	}

	// Encryption
	if db.EncryptData != nil {
		if db.EncryptData.KeyId != nil {
			details["Encryption Key"] = *db.EncryptData.KeyId
		} else {
			details["Encryption Key"] = "Oracle-managed key"
		}
		orderedKeys = append(orderedKeys, "Encryption Key")
	}

	// Security certificates
	if db.SecureConnections != nil {
		if db.SecureConnections.CertificateGenerationType != "" {
			details["Security Certificate"] = string(db.SecureConnections.CertificateGenerationType)
			orderedKeys = append(orderedKeys, "Security Certificate")
		}
	}

	// Read endpoint
	if db.ReadEndpoint != nil {
		if db.ReadEndpoint.ReadEndpointIpAddress != nil {
			details["Read Endpoint IP"] = *db.ReadEndpoint.ReadEndpointIpAddress
			orderedKeys = append(orderedKeys, "Read Endpoint IP")
		} else {
			details["Read Endpoint"] = "Disabled"
			orderedKeys = append(orderedKeys, "Read Endpoint")
		}
	}

	// Database Management
	if db.DatabaseManagement != "" {
		details["Database Management"] = db.DatabaseManagement
		orderedKeys = append(orderedKeys, "Database Management")
	}

	// Customer contacts
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

	// Lifecycle details (useful for troubleshooting)
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
		"MySQL.Free": {ecpu: 2, memory: 8},      // Free tier: 2 ECPUs, 8 GB
		"MySQL.2":    {ecpu: 2, memory: 16},     // MySQL.2: 2 ECPUs, 16 GB
		"MySQL.4":    {ecpu: 4, memory: 32},     // MySQL.4: 4 ECPUs, 32 GB
		"MySQL.8":    {ecpu: 8, memory: 64},     // MySQL.8: 8 ECPUs, 64 GB
		"MySQL.16":   {ecpu: 16, memory: 128},   // MySQL.16: 16 ECPUs, 128 GB
		"MySQL.32":   {ecpu: 32, memory: 256},   // MySQL.32: 32 ECPUs, 256 GB
		"MySQL.48":   {ecpu: 48, memory: 384},   // MySQL.48: 48 ECPUs, 384 GB
		"MySQL.64":   {ecpu: 64, memory: 512},   // MySQL.64: 64 ECPUs, 512 GB
		"MySQL.256":  {ecpu: 256, memory: 1024}, // MySQL.256: 256 ECPUs, 1024 GB (2 x 50 Gbps)
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
