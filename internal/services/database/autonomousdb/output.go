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
func PrintAutonomousDbInfo(databases []domain.AutonomousDatabase, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
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

		// Connection (existing behavior) + includes TP/TPURGENT if present
		conn := map[string]string{
			"Private IP":       db.PrivateEndpointIp,
			"Private Endpoint": db.PrivateEndpoint,
			"High":             db.ConnectionStrings["HIGH"],
			"Medium":           db.ConnectionStrings["MEDIUM"],
			"Low":              db.ConnectionStrings["LOW"],
			"TP":               db.ConnectionStrings["TP"],
			"TPURGENT":         db.ConnectionStrings["TPURGENT"],
		}
		connKeys := []string{"Private IP", "Private Endpoint", "High", "Medium", "Low", "TP", "TPURGENT"}
		p.PrintKeyValues(title, conn, connKeys)

		// State section
		state := map[string]string{
			"Lifecycle State":   db.LifecycleState,
			"Lifecycle Details": db.LifecycleDetails,
			"DB Version":        db.DbVersion,
			"Workload":          db.DbWorkload,
			"License Model":     db.LicenseModel,
		}
		if db.TimeCreated != nil {
			state["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
		}
		stateKeys := []string{"Lifecycle State", "Lifecycle Details", "DB Version", "Workload", "License Model", "Time Created"}
		p.PrintKeyValues(title+" – State", state, stateKeys)

		// Network section (prefer names over OCIDs when available)
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
		net := map[string]string{
			"Private Endpoint Label": db.PrivateEndpointLabel,
			"Access Type":            accessType,
			"Subnet":                 subnetVal,
			"VCN":                    vcnVal,
			"NSGs":                   nsgVal,
			"mTLS Required":          boolToString(db.IsMtlsRequired),
		}
		if len(db.WhitelistedIps) > 0 {
			net["Whitelisted IPs"] = fmt.Sprintf("%v", db.WhitelistedIps)
		}
		netKeys := []string{"Private Endpoint Label", "Access Type", "Subnet", "VCN", "NSGs", "Whitelisted IPs", "mTLS Required"}
		p.PrintKeyValues(title+" – Network", net, netKeys)

		// Capacity section: show ECPUs when ComputeModel is ECPU, otherwise OCPUs/CPU Cores
		cap := map[string]string{}
		var capKeys []string
		// Determine storage value: prefer TBs, else GBs if present
		storage := ""
		if s := intToString(db.DataStorageSizeInTBs); s != "" {
			storage = s + " TB"
		} else if g := intToString(db.DataStorageSizeInGBs); g != "" {
			storage = g + " GB"
		}
		if db.ComputeModel == "ECPU" || db.EcpuCount != nil {
			cap["Compute Model"] = db.ComputeModel
			cap["ECPUs"] = floatToString(db.EcpuCount)
			cap["Storage"] = storage
			cap["Auto Scaling"] = boolToString(db.IsAutoScalingEnabled)
			capKeys = []string{"Compute Model", "ECPUs", "Storage", "Auto Scaling"}
		} else {
			cap["Compute Model"] = db.ComputeModel
			cap["OCPUs"] = floatToString(db.OcpuCount)
			cap["CPU Cores"] = intToString(db.CpuCoreCount)
			cap["Storage"] = storage
			cap["Auto Scaling"] = boolToString(db.IsAutoScalingEnabled)
			capKeys = []string{"Compute Model", "OCPUs", "CPU Cores", "Storage", "Auto Scaling"}
		}
		p.PrintKeyValues(title+" – Capacity", cap, capKeys)
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
