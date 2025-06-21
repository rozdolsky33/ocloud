package oke

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

func PrintOKEInfo(clusters []Cluster, appCtx *app.ApplicationContext, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return marshalOKEToJSON(p, clusters)
	}

	// Handle the case where no clusters are found.
	if len(clusters) == 0 {
		fmt.Fprintln(appCtx.Stdout, "No clusters found.")
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

		// Create the colored title using components from the app context.
		coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
		coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
		coloredCluster := text.Colors{text.FgBlue}.Sprint(cluster.Name)
		title := fmt.Sprintf("%s: %s: %s",
			coloredTenancy,
			coloredCompartment,
			coloredCluster)

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
					coloredCluster)

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

// marshalOKEToJSON marshals clusters to JSON format.
// It accepts a printer and returns an error.
func marshalOKEToJSON(p *printer.Printer, clusters []Cluster) error {
	response := JSONResponse{
		Clusters: clusters,
	}
	// Use the printer's method to marshal. It will write to the correct output.
	return p.MarshalToJSON(response)
}
