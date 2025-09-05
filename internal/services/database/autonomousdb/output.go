package autonomousdb

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintAutonomousDbInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintAutonomousDbInfo(databases []domain.AutonomousDatabase, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		if len(databases) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
		return util.MarshalDataToJSONResponse[domain.AutonomousDatabase](p, databases, pagination)
	}

	if util.ValidateAndReportEmpty(databases, pagination, appCtx.Stdout) {
		return nil
	}
	// Print each database as a set of key-value sections with a colored title.
	for _, db := range databases {
		title := util.FormatColoredTitle(appCtx, db.Name)

		// Prepare name-preferred fields for subnet/VCN and NSGs
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
		// Storage: prefer TBs than GBs
		storage := ""
		if s := intToString(db.DataStorageSizeInTBs); s != "" {
			storage = s + " TB"
		} else if g := intToString(db.DataStorageSizeInGBs); g != "" {
			storage = g + " GB"
		}
		// CPU label/value based on compute model
		cpuKey := "OCPUs"
		cpuVal := floatToString(db.OcpuCount)
		if db.ComputeModel == "ECPU" || db.EcpuCount != nil {
			cpuKey = "ECPUs"
			cpuVal = floatToString(db.EcpuCount)
		}

		// Summary view (first glance) when showAll is false
		if !showAll {
			summary := map[string]string{
				"Lifecycle State":  db.LifecycleState,
				"DB Version":       db.DbVersion,
				"Workload":         db.DbWorkload,
				"Compute Model":    db.ComputeModel,
				cpuKey:             cpuVal,
				"Storage":          storage,
				"Private IP":       db.PrivateEndpointIp,
				"Private Endpoint": db.PrivateEndpoint,
				"Subnet":           subnetVal,
				"VCN":              vcnVal,
				"NSGs":             nsgVal,
			}
			if db.TimeCreated != nil {
				summary["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
			}
			// Ordered keys for summary
			ordered := []string{"Lifecycle State", "DB Version", "Workload", "Compute Model", cpuKey, "Storage", "Private IP", "Private Endpoint", "Subnet", "VCN", "NSGs", "Time Created"}
			p.PrintKeyValues(title, summary, ordered)
			continue
		}

		// Detailed view (showAll is true)
		details := make(map[string]string)
		orderedKeys := []string{}

		// General Section
		details["Lifecycle State"] = db.LifecycleState
		details["DB Version"] = db.DbVersion
		details["Workload"] = db.DbWorkload
		details["License Model"] = db.LicenseModel
		if db.TimeCreated != nil {
			details["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
		}
		orderedKeys = append(orderedKeys, "Lifecycle State", "Lifecycle Details", "DB Version", "Workload", "License Model", "Time Created")

		// Capacity Section
		details["Compute Model"] = db.ComputeModel
		if db.ComputeModel == "ECPU" || db.EcpuCount != nil {
			details["ECPUs"] = floatToString(db.EcpuCount)
			orderedKeys = append(orderedKeys, "Compute Model", "ECPUs")
		} else {
			details["OCPUs"] = floatToString(db.OcpuCount)
			details["CPU Cores"] = intToString(db.CpuCoreCount)
			orderedKeys = append(orderedKeys, "Compute Model", "OCPUs", "CPU Cores")
		}
		details["Storage"] = storage
		details["Auto Scaling"] = boolToString(db.IsAutoScalingEnabled)
		orderedKeys = append(orderedKeys, "Storage", "Auto Scaling")

		// Network Section
		accessType := ""
		if db.IsPubliclyAccessible != nil {
			if *db.IsPubliclyAccessible {
				accessType = "Public"
			} else if db.PrivateEndpoint != "" {
				accessType = "Virtual cloud network"
			}
		} else if db.PrivateEndpoint != "" {
			// Infer VCN when we have a private endpoint but no explicit public flag
			accessType = "Virtual cloud network"
		}
		details["Access Type"] = accessType
		details["Private IP"] = db.PrivateEndpointIp
		details["Private Endpoint"] = db.PrivateEndpoint
		details["Subnet"] = subnetVal
		details["VCN"] = vcnVal
		details["NSGs"] = nsgVal
		details["mTLS Required"] = boolToString(db.IsMtlsRequired)
		if len(db.WhitelistedIps) > 0 {
			details["Whitelisted IPs"] = fmt.Sprintf("%v", db.WhitelistedIps)
		}
		orderedKeys = append(orderedKeys, "Access Type", "Private IP", "Private Endpoint", "Subnet", "VCN", "NSGs", "mTLS Required", "Whitelisted IPs")

		// Connection Strings Section
		details["High"] = db.ConnectionStrings["HIGH"]
		details["Medium"] = db.ConnectionStrings["MEDIUM"]
		details["Low"] = db.ConnectionStrings["LOW"]
		details["TP"] = db.ConnectionStrings["TP"]
		details["TPURGENT"] = db.ConnectionStrings["TPURGENT"]
		orderedKeys = append(orderedKeys, "High", "Medium", "Low", "TP", "TPURGENT")

		p.PrintKeyValues(title, details, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// helpers for pointer formatting
func boolToString(v *bool) string {
	if v == nil {
		return ""
	}
	if *v {
		return "true"
	}
	return "false"
}

func intToString(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}

func floatToString(v *float32) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *v)
}
