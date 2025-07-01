package oke

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintOKETable groups cluster metadata and node‑pool details into **one**
// responsive table per cluster. The first row summarizes the cluster; the
// following rows list each node pool.
func PrintOKETable(clusters []Cluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

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

// PrintOKEInfo prints a detailed, troubleshooting‑oriented view of OKE clusters
// and their node pools.  Each cluster is rendered as:
//  1. A summary key/value block containing operationally‑relevant metadata.
//  2. A responsive table listing all node pools with the most useful columns
//     for SRE triage.
//
// When --JSON is requested, the function defers to util.MarshalDataToJSONResponse.
func PrintOKEInfo(clusters []Cluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[Cluster](p, clusters, pagination)
	}

	if util.ValidateAndReportEmpty(clusters, pagination, appCtx.Stdout) {
		return nil
	}

	// Sort clusters by name for deterministic output.
	sort.Slice(clusters, func(i, j int) bool {
		return strings.ToLower(clusters[i].Name) < strings.ToLower(clusters[j].Name)
	})

	for _, c := range clusters {
		summary := map[string]string{
			"ID":               c.ID,
			"Name":             c.Name,
			"K8s Version":      c.Version,
			"Created":          c.CreatedAt,
			"State":            string(c.State),
			"Private Endpoint": c.PrivateEndpoint,
			"Node Pools":       fmt.Sprintf("%d", len(c.NodePools)),
		}

		order := []string{"ID", "Name", "K8s Version", "Created", "State", "Private Endpoint", "Node Pools"}

		title := util.FormatColoredTitle(appCtx, fmt.Sprintf("Cluster: %s", c.Name))
		p.PrintKeyValues(title, summary, order)
		fmt.Fprintln(appCtx.Stdout) // spacer

		//-----------------------------------------------------------------
		// Node pool details (table)
		//-----------------------------------------------------------------
		if len(c.NodePools) > 0 {
			headers := []string{"Node Pool", "Version", "Shape", "OCPUs", "Mem(GB)", "Node Cnt", "State"}
			rows := make([][]string, len(c.NodePools))

			for i, np := range c.NodePools {
				rows[i] = []string{
					np.Name,
					np.Version,
					np.NodeShape,
					np.Ocpus,
					np.MemoryGB,
					fmt.Sprintf("%d", np.NodeCount),
					string(np.State),
				}
			}

			tableTitle := util.FormatColoredTitle(appCtx, "Node Pools")
			p.PrintTable(tableTitle, headers, rows)
			fmt.Fprintln(appCtx.Stdout) // spacer between clusters
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
