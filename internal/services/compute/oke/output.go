package oke

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintOKEInfo prints Oracle Kubernetes Engine clusters and their node pools in
// a compact, grouped layout.  Each cluster is shown as a summary key/value
// block followed by a single table that lists all its node pools.
func PrintOKEInfo(clusters []Cluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	// Handle pagination meta‑data if present.
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// JSON output?  Hand off early.
	if useJSON {
		return util.MarshalDataToJSONResponse[Cluster](p, clusters, pagination)
	}

	// Empty check.
	if util.ValidateAndReportEmpty(clusters, pagination, appCtx.Stdout) {
		return nil
	}

	for _, cluster := range clusters {
		// -------------------- Cluster summary block --------------------
		clusterSummary := map[string]string{
			"ID":               cluster.ID,
			"Name":             cluster.Name,
			"Version":          cluster.Version,
			"Created":          cluster.CreatedAt,
			"State":            string(cluster.State),
			"Private Endpoint": cluster.PrivateEndpoint,
		}

		summaryKeys := []string{"ID", "Name", "Version", "Created", "State", "Private Endpoint"}

		summaryTitle := util.FormatColoredTitle(appCtx, fmt.Sprintf("Cluster: %s", cluster.Name))
		p.PrintKeyValues(summaryTitle, clusterSummary, summaryKeys)
		fmt.Fprintln(appCtx.Stdout) // spacer

		// ----------------------- Node pool table -----------------------
		if len(cluster.NodePools) > 0 {
			headers := []string{"Node Pool", "Version", "Shape", "Node Count", "State"}
			rows := make([][]string, len(cluster.NodePools))

			for i, np := range cluster.NodePools {
				rows[i] = []string{
					np.Name,
					np.Version,
					np.NodeShape,
					fmt.Sprintf("%d", np.NodeCount),
					string(np.State),
				}
			}

			tableTitle := util.FormatColoredTitle(appCtx, fmt.Sprintf("Node Pools (%d)", len(cluster.NodePools)))
			p.PrintTable(tableTitle, headers, rows)
			fmt.Fprintln(appCtx.Stdout) // spacer between clusters
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
