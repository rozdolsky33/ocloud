package cacheclusterdb

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintCacheClusterInfo prints a single HeatWave Cache Cluster.
func PrintCacheClusterInfo(cluster *database.CacheCluster, appCtx *app.ApplicationContext, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(cluster)
	}

	return printOneCacheCluster(p, appCtx, cluster, showAll)
}

// PrintCacheClustersInfo prints a list of HeatWave Cache Clusters.
func PrintCacheClustersInfo(clusters []database.CacheCluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		if len(clusters) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
		return util.MarshalDataToJSONResponse[database.CacheCluster](p, clusters, pagination)
	}

	if util.ValidateAndReportEmpty(clusters, pagination, appCtx.Stdout) {
		return nil
	}

	for _, cluster := range clusters {
		if err := printOneCacheCluster(p, appCtx, &cluster, showAll); err != nil {
			return err
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func printOneCacheCluster(p *printer.Printer, appCtx *app.ApplicationContext, cluster *database.CacheCluster, showAll bool) error {
	title := util.FormatColoredTitle(appCtx, cluster.DisplayName)

	subnetVal := cluster.SubnetId
	if cluster.SubnetName != "" {
		subnetVal = cluster.SubnetName
	}
	vcnVal := cluster.VcnID
	if cluster.VcnName != "" {
		vcnVal = cluster.VcnName
	}

	// NSG information
	nsgVal := ""
	if len(cluster.NsgNames) > 0 {
		nsgVal = fmt.Sprintf("%v", cluster.NsgNames)
	} else if len(cluster.NsgIds) > 0 {
		nsgVal = fmt.Sprintf("%v", cluster.NsgIds)
	} else {
		nsgVal = "None"
	}

	// Cluster mode info
	clusterModeInfo := cluster.ClusterMode
	if cluster.ClusterMode == "SHARDED" && cluster.ShardCount > 0 {
		clusterModeInfo = fmt.Sprintf("%s (%d shards)", cluster.ClusterMode, cluster.ShardCount)
	}

	// Node information
	nodeInfo := fmt.Sprintf("%d nodes", cluster.NodeCount)
	if cluster.NodeMemoryInGBs > 0 {
		nodeInfo = fmt.Sprintf("%d nodes × %.0fGB", cluster.NodeCount, cluster.NodeMemoryInGBs)
	}

	if !showAll {
		// Summary view - Essential operational info
		summary := map[string]string{
			"Lifecycle State":   cluster.LifecycleState,
			"Software Version":  cluster.SoftwareVersion,
			"Cluster Mode":      clusterModeInfo,
			"Nodes":             nodeInfo,
			"Primary Endpoint":  cluster.PrimaryFqdn,
			"Replicas Endpoint": cluster.ReplicasFqdn,
			"Subnet":            subnetVal,
			"VCN":               vcnVal,
		}

		if cluster.TimeCreated != nil {
			summary["Time Created"] = cluster.TimeCreated.Format("2006-01-02 15:04:05")
		}

		ordered := []string{
			"Lifecycle State", "Software Version", "Cluster Mode", "Nodes",
			"Primary Endpoint", "Replicas Endpoint", "Subnet", "VCN", "Time Created",
		}
		p.PrintKeyValues(title, summary, ordered)
		return nil
	}

	// Detailed view
	details := make(map[string]string)
	orderedKeys := []string{}

	// General
	details["Lifecycle State"] = cluster.LifecycleState
	if cluster.LifecycleDetails != "" {
		details["Lifecycle Details"] = cluster.LifecycleDetails
	}
	if cluster.TimeCreated != nil {
		details["Time Created"] = cluster.TimeCreated.Format("2006-01-02 15:04:05")
	}
	if cluster.TimeUpdated != nil {
		details["Time Updated"] = cluster.TimeUpdated.Format("2006-01-02 15:04:05")
	}

	generalKeys := []string{"Lifecycle State"}
	if cluster.LifecycleDetails != "" {
		generalKeys = append(generalKeys, "Lifecycle Details")
	}
	generalKeys = append(generalKeys, "Time Created", "Time Updated")
	orderedKeys = append(orderedKeys, generalKeys...)

	// Cluster Configuration
	details["Software Version"] = cluster.SoftwareVersion
	details["Cluster Mode"] = clusterModeInfo
	details["Nodes"] = nodeInfo
	if cluster.ConfigSetId != "" {
		details["Config Set ID"] = cluster.ConfigSetId
	}
	orderedKeys = append(orderedKeys, "Software Version", "Cluster Mode", "Nodes")
	if cluster.ConfigSetId != "" {
		orderedKeys = append(orderedKeys, "Config Set ID")
	}

	// Endpoints
	details["Primary FQDN"] = cluster.PrimaryFqdn
	details["Primary IP"] = cluster.PrimaryEndpointIpAddress
	details["Replicas FQDN"] = cluster.ReplicasFqdn
	details["Replicas IP"] = cluster.ReplicasEndpointIpAddress
	orderedKeys = append(orderedKeys, "Primary FQDN", "Primary IP", "Replicas FQDN", "Replicas IP")

	if cluster.DiscoveryFqdn != "" {
		details["Discovery FQDN"] = cluster.DiscoveryFqdn
		details["Discovery IP"] = cluster.DiscoveryEndpointIpAddress
		orderedKeys = append(orderedKeys, "Discovery FQDN", "Discovery IP")
	}

	// Individual Nodes
	if len(cluster.Nodes) > 0 {
		for i, node := range cluster.Nodes {
			if node.DisplayName != nil {
				nodeKey := fmt.Sprintf("Node %d Name", i+1)
				details[nodeKey] = *node.DisplayName
				orderedKeys = append(orderedKeys, nodeKey)

				if node.PrivateEndpointFqdn != nil {
					endpointKey := fmt.Sprintf("Node %d Endpoint", i+1)
					details[endpointKey] = *node.PrivateEndpointFqdn
					orderedKeys = append(orderedKeys, endpointKey)
				}
			}
		}
	}

	// Network
	details["Subnet"] = subnetVal
	details["VCN"] = vcnVal
	if nsgVal != "None" {
		details["NSGs"] = nsgVal
		orderedKeys = append(orderedKeys, "Subnet", "VCN", "NSGs")
	} else {
		orderedKeys = append(orderedKeys, "Subnet", "VCN")
	}

	p.PrintKeyValues(title, details, orderedKeys)
	return nil
}
