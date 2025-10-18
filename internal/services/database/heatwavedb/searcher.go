package heatwavedb

import (
	"strconv"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SearchableHeatWaveDatabase adapts HeatWaveDatabase to the search.Indexable interface.
type SearchableHeatWaveDatabase struct {
	database.HeatWaveDatabase
}

// ToIndexable converts a HeatWaveDatabase to a map of searchable fields.
func (s SearchableHeatWaveDatabase) ToIndexable() map[string]any {
	tagsKV, _ := util.FlattenTags(s.FreeformTags, s.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(s.FreeformTags, s.DefinedTags)

	// join slices safely
	join := func(items []string) string {
		return strings.ToLower(strings.Join(items, ","))
	}

	// Format storage size
	var storageGB string
	if s.DataStorage != nil && s.DataStorage.DataStorageSizeInGBs != nil {
		storageGB = strconv.Itoa(*s.DataStorage.DataStorageSizeInGBs)
	} else if s.DataStorageSizeInGBs != nil {
		storageGB = strconv.Itoa(*s.DataStorageSizeInGBs)
	}

	// Format HeatWave cluster size
	var clusterSize string
	if s.HeatWaveCluster != nil && s.HeatWaveCluster.ClusterSize != nil {
		clusterSize = strconv.Itoa(*s.HeatWaveCluster.ClusterSize)
	}

	return map[string]any{
		"Name":               strings.ToLower(s.DisplayName),
		"OCID":               strings.ToLower(s.ID),
		"State":              strings.ToLower(s.LifecycleState),
		"Description":        strings.ToLower(s.Description),
		"MysqlVersion":       strings.ToLower(s.MysqlVersion),
		"ShapeName":          strings.ToLower(s.ShapeName),
		"StorageGB":          strings.ToLower(storageGB),
		"DatabaseMode":       strings.ToLower(s.DatabaseMode),
		"AccessMode":         strings.ToLower(s.AccessMode),
		"VcnID":              strings.ToLower(s.VcnID),
		"VcnName":            strings.ToLower(s.VcnName),
		"SubnetId":           strings.ToLower(s.SubnetId),
		"SubnetName":         strings.ToLower(s.SubnetName),
		"HostnameLabel":      strings.ToLower(s.HostnameLabel),
		"IpAddress":          strings.ToLower(s.IpAddress),
		"NsgNames":           join(s.NsgNames),
		"NsgIds":             join(s.NsgIds),
		"ClusterSize":        strings.ToLower(clusterSize),
		"AvailabilityDomain": strings.ToLower(s.AvailabilityDomain),
		"FaultDomain":        strings.ToLower(s.FaultDomain),
		"CrashRecovery":      strings.ToLower(s.CrashRecovery),
		"TagsKV":             strings.ToLower(tagsKV),
		"TagsVal":            strings.ToLower(tagsVal),
	}
}

// GetSearchableFields returns the list of fields to be indexed for HeatWave Databases.
func GetSearchableFields() []string {
	return []string{
		"Name", "OCID", "State", "Description", "MysqlVersion", "ShapeName",
		"StorageGB", "DatabaseMode", "AccessMode",
		"VcnID", "VcnName", "SubnetId", "SubnetName",
		"HostnameLabel", "IpAddress",
		"NsgNames", "NsgIds",
		"ClusterSize", "AvailabilityDomain", "FaultDomain", "CrashRecovery",
		"TagsKV", "TagsVal",
	}
}

// GetBoostedFields returns the list of fields to be boosted in the search.
func GetBoostedFields() []string {
	return []string{"Name", "OCID", "VcnName", "SubnetName", "IpAddress"}
}

// ToSearchableHeatWaveDbs converts a slice of HeatWaveDatabase to a slice of search.Indexable.
func ToSearchableHeatWaveDbs(dbs []database.HeatWaveDatabase) []search.Indexable {
	searchable := make([]search.Indexable, len(dbs))
	for i, db := range dbs {
		searchable[i] = SearchableHeatWaveDatabase{db}
	}
	return searchable
}
