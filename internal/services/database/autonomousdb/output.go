package autonomousdb

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintAutonomousDbInfo prints a single Autonomous DB.
// - useJSON: if true, prints the single DB as JSON (no pagination envelope)
// - showAll: if true, prints the detailed view; otherwise, prints the summary view
func PrintAutonomousDbInfo(db *domain.AutonomousDatabase, appCtx *app.ApplicationContext, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(db)
	}

	return printOneAutonomousDb(p, appCtx, db, showAll)
}

// PrintAutonomousDbsInfo prints a list of Autonomous DBs.
// - pagination: optional, will be adjusted and logged if provided
// - useJSON: if true, prints databases with util.MarshalDataToJSONResponse
// - showAll: if true, prints detailed view; otherwise summary view
func PrintAutonomousDbsInfo(databases []domain.AutonomousDatabase, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		if len(databases) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
		return util.MarshalDataToJSONResponse[domain.AutonomousDatabase](p, databases, pagination)
	}

	if util.ValidateAndReportEmpty(databases, pagination, appCtx.Stdout) {
		return nil
	}

	for _, db := range databases {
		if err := printOneAutonomousDb(p, appCtx, &db, showAll); err != nil {
			return err
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func printOneAutonomousDb(p *printer.Printer, appCtx *app.ApplicationContext, db *domain.AutonomousDatabase, showAll bool) error {
	title := util.FormatColoredTitle(appCtx, db.Name)

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

	// Storage: prefer TBs over GBs
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

	if !showAll {
		// Summary view
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
		ordered := []string{
			"Lifecycle State", "DB Version", "Workload", "Compute Model",
			cpuKey, "Storage", "Private IP", "Private Endpoint", "Subnet", "VCN", "NSGs", "Time Created",
		}
		p.PrintKeyValues(title, summary, ordered)
		return nil
	}

	// Detailed view
	details := make(map[string]string)
	orderedKeys := []string{}

	// General
	details["Lifecycle State"] = db.LifecycleState
	details["DB Version"] = db.DbVersion
	details["Workload"] = db.DbWorkload
	details["License Model"] = db.LicenseModel
	if db.TimeCreated != nil {
		details["Time Created"] = db.TimeCreated.Format("2006-01-02 15:04:05")
	}
	orderedKeys = append(orderedKeys, "Lifecycle State", "Lifecycle Details", "DB Version", "Workload", "License Model", "Time Created")

	// Capacity
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

	// Network
	accessType := ""
	if db.IsPubliclyAccessible != nil {
		if *db.IsPubliclyAccessible {
			accessType = "Public"
		} else if db.PrivateEndpoint != "" {
			accessType = "Virtual cloud network"
		}
	} else if db.PrivateEndpoint != "" {
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

	// Connection Strings
	if details["High"] = db.ConnectionStrings["HIGH"]; details["High"] != "" {
		orderedKeys = append(orderedKeys, "High")
	}
	if details["Medium"] = db.ConnectionStrings["MEDIUM"]; details["Medium"] != "" {
		orderedKeys = append(orderedKeys, "Medium")
	}
	if details["Low"] = db.ConnectionStrings["LOW"]; details["Low"] != "" {
		orderedKeys = append(orderedKeys, "Low")
	}
	if details["TP"] = db.ConnectionStrings["TP"]; details["TP"] != "" {
		orderedKeys = append(orderedKeys, "TP")
	}
	if details["TPURGENT"] = db.ConnectionStrings["TPURGENT"]; details["TPURGENT"] != "" {
		orderedKeys = append(orderedKeys, "TPURGENT")
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
