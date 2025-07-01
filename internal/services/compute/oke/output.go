package oke

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintOKEInfo groups cluster metadata and node‑pool details into **one**
// responsive table per cluster. The first row summarizes the cluster; the
// following rows list each node pool.
func PrintOKEInfo(clusters []Cluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// JSON pathway --------------------------------------------------------
	if useJSON {
		return util.MarshalDataToJSONResponse[Cluster](p, clusters, pagination)
	}

	if util.ValidateAndReportEmpty(clusters, pagination, appCtx.Stdout) {
		return nil
	}

	for _, c := range clusters {
		// Header defines unified columns for both cluster + node‑pool rows.
		headers := []string{
			"Name",           // cluster or node‑pool name
			"Type",           // "Cluster" or "NodePool"
			"Version",        // k8s version
			"Shape/Endpoint", // node shape or cluster endpoint
			"Count/Created",  // node count or created timestamp
			"State",          // lifecycle state
		}

		// Build rows — first row is the cluster summary.
		rows := [][]string{
			{
				c.Name,
				"Cluster",
				c.Version,
				c.PrivateEndpoint,
				c.CreatedAt,
				string(c.State),
			},
		}

		// Append node‑pool rows.
		for _, np := range c.NodePools {
			rows = append(rows, []string{
				np.Name,
				"NodePool",
				np.Version,
				np.NodeShape,
				fmt.Sprintf("%d", np.NodeCount),
				string(np.State),
			})
		}

		title := util.FormatColoredTitle(appCtx, fmt.Sprintf("Cluster: %s (%d node pools)", c.Name, len(c.NodePools)))
		p.PrintTable(title, headers, rows)
		fmt.Fprintln(appCtx.Stdout) // spacer between clusters
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
