package oke

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintOKEInfo(clusters []Cluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return util.MarshalDataToJSONResponse[Cluster](p, clusters, pagination)
	}

	if util.ValidateAndReportEmpty(clusters, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each cluster as a separate key-value table with a colored title.
	for _, cluster := range clusters {
		// Create cluster data map
		clusterData := map[string]string{
			"ID":               cluster.ID,
			"Name":             cluster.Name,
			"Version":          cluster.Version,
			"Created":          cluster.CreatedAt,
			"State":            string(cluster.State),
			"Private Endpoint": cluster.PrivateEndpoint,
			"Node Pools Count": fmt.Sprintf("%d", len(cluster.NodePools)),
		}

		// Define ordered keys
		orderedKeys := []string{
			"ID", "Name", "Version", "Created", "State", "Private Endpoint", "Node Pools Count",
		}

		title := util.FormatColoredTitle(appCtx, cluster.Name)

		// Call the printer method to render the key-value table for this cluster.
		p.PrintKeyValues(title, clusterData, orderedKeys)

		// Print node pool details if there are any
		if len(cluster.NodePools) > 0 {
			fmt.Fprintln(appCtx.Stdout, "\nNode Pools:", len(cluster.NodePools))

			// Print each cluster as a separate key-value table with a colored title.
			for _, node := range cluster.NodePools {
				// Create cluster data map
				nodePoolData := map[string]string{
					"ID":         node.ID,
					"Version":    node.Version,
					"Shape":      node.NodeShape,
					"Node Count": fmt.Sprintf("%d", node.NodeCount),
					"State":      string(node.State),
				}

				// Define ordered keys
				nodePoolKeys := []string{
					"ID", "Version", "Shape", "Node Count", "State"}

				// Create the colored title using components from the app context.
				coloredNodePool := text.Colors{text.FgMagenta}.Sprint("Node Pool")
				nodeTitle := fmt.Sprintf("%s: %s: %s",
					coloredNodePool,
					node.Name,
					cluster.Name)

				// Call the printer method to render the key-value table for this cluster.
				p.PrintKeyValues(nodeTitle, nodePoolData, nodePoolKeys)

				// Add a separator between clusters
				fmt.Fprintln(appCtx.Stdout, "")
			}

			return nil
		}
	}
	return nil
}
