package autonomousdb

import (
	"strconv"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableAutonomousDatabase adapts AutonomousDatabase to the search.Indexable interface.
type SearchableAutonomousDatabase struct {
	database.AutonomousDatabase
}

// ToIndexable converts an AutonomousDatabase to a map of searchable fields.
func (s SearchableAutonomousDatabase) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	// join slices safely
	join := func(items []string) string {
		return strings.ToLower(strings.Join(items, ","))
	}

	return map[string]any{
		"Name":               strings.ToLower(s.Name),
		"OCID":               strings.ToLower(s.ID),
		"State":              strings.ToLower(s.LifecycleState),
		"DbVersion":          strings.ToLower(s.DbVersion),
		"Workload":           strings.ToLower(s.DbWorkload),
		"LicenseModel":       strings.ToLower(s.LicenseModel),
		"ComputeModel":       strings.ToLower(s.ComputeModel),
		"OcpuCount":          strings.ToLower(formatFloat32Ptr(s.OcpuCount)),
		"EcpuCount":          strings.ToLower(formatFloat32Ptr(s.EcpuCount)),
		"CpuCoreCount":       strings.ToLower(formatIntPtr(s.CpuCoreCount)),
		"StorageTB":          strings.ToLower(formatIntPtr(s.DataStorageSizeInTBs)),
		"StorageGB":          strings.ToLower(formatIntPtr(s.DataStorageSizeInGBs)),
		"VcnID":              strings.ToLower(s.VcnID),
		"VcnName":            strings.ToLower(s.VcnName),
		"SubnetId":           strings.ToLower(s.SubnetId),
		"SubnetName":         strings.ToLower(s.SubnetName),
		"PrivateEndpoint":    strings.ToLower(s.PrivateEndpoint),
		"PrivateEndpointIp":  strings.ToLower(s.PrivateEndpointIp),
		"PrivateEndpointLbl": strings.ToLower(s.PrivateEndpointLabel),
		"WhitelistedIps":     join(s.WhitelistedIps),
		"NsgNames":           join(s.NsgNames),
		"NsgIds":             join(s.NsgIds),
		"TagsKV":             strings.ToLower(tagsKV),
		"TagsVal":            strings.ToLower(tagsVal),
	}
}

// formatFloat32Ptr returns string value of *float32 or empty if nil.
func formatFloat32Ptr(p *float32) string {
	if p == nil {
		return ""
	}
	return strconv.FormatFloat(float64(*p), 'f', -1, 64)
}

// formatIntPtr returns string value of *int or empty if nil.
func formatIntPtr(p *int) string {
	if p == nil {
		return ""
	}
	return strconv.Itoa(*p)
}

// GetSearchableFields returns the list of fields to be indexed for Autonomous Databases.
func GetSearchableFields() []string {
	return []string{
		"Name", "OCID", "State", "DbVersion", "Workload", "LicenseModel",
		"ComputeModel", "OcpuCount", "EcpuCount", "CpuCoreCount", "StorageTB", "StorageGB",
		"VcnID", "VcnName", "SubnetId", "SubnetName",
		"PrivateEndpoint", "PrivateEndpointIp", "PrivateEndpointLbl",
		"WhitelistedIps", "NsgNames", "NsgIds",
		"TagsKV", "TagsVal",
	}
}

// GetBoostedFields returns the list of fields to be boosted in the search.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "VcnName", "SubnetName"}
}

// ToSearchableAutonomousDBs converts a slice of AutonomousDatabase to a slice of search.Indexable.
func ToSearchableAutonomousDBs(dbs []AutonomousDatabase) []search.Indexable {
	searchable := make([]search.Indexable, len(dbs))
	for i, db := range dbs {
		searchable[i] = SearchableAutonomousDatabase{db}
	}
	return searchable
}
