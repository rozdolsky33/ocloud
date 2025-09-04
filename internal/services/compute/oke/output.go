package oke

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintOKETable groups cluster metadata and node-pool details into one table per cluster.
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
		headers := []string{"Name", "Type", "Version", "Shape/Endpoint", "Count/Created", "State"}

		rows := [][]string{
			{
				c.DisplayName,
				"Cluster",
				c.KubernetesVersion,
				c.PrivateEndpoint,
				c.TimeCreated.Format("2006-01-02"),
				c.State,
			},
		}

		for _, np := range c.NodePools {
			rows = append(rows, []string{
				np.DisplayName,
				"NodePool",
				np.KubernetesVersion,
				np.NodeShape,
				fmt.Sprintf("%d", np.NodeCount),
				"", // State for node pools is not in the domain model yet
			})
		}

		title := util.FormatColoredTitle(appCtx, fmt.Sprintf("Cluster: %s (%d node pools)", c.DisplayName, len(c.NodePools)))
		p.PrintTable(title, headers, rows)
		fmt.Fprintln(appCtx.Stdout)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// PrintOKEsInfo displays instances in a formatted table or JSON format.
func PrintOKEsInfo(clusters []Cluster, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
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

	sort.Slice(clusters, func(i, j int) bool {
		return strings.ToLower(clusters[i].DisplayName) < strings.ToLower(clusters[j].DisplayName)
	})

	for _, c := range clusters {
		renderCluster(p, appCtx, c)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// PrintOKEInfo prints a detailed view of a cluster.
func PrintOKEInfo(appCtx *app.ApplicationContext, c *Cluster, useJSON bool) error {
	if c == nil {
		return fmt.Errorf("nil cluster")
	}
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return util.MarshalDataToJSONResponse[Cluster](p, []Cluster{*c}, nil)
	}

	renderCluster(p, appCtx, *c)
	return nil
}

// renderCluster renders a cluster in a formatted table or JSON format.
func renderCluster(p *printer.Printer, appCtx *app.ApplicationContext, c Cluster) {
	created := "-"
	if !c.TimeCreated.IsZero() {
		created = c.TimeCreated.Format("2006-01-02 15:04:05")
	}

	summary := map[string]string{
		"ID":               c.OCID,
		"Name":             c.DisplayName,
		"K8s Version":      c.KubernetesVersion,
		"Created":          created,
		"State":            c.State,
		"Private Endpoint": c.PrivateEndpoint,
		"Node Pools":       fmt.Sprintf("%d", len(c.NodePools)),
	}
	order := []string{"ID", "Name", "K8s Version", "Created", "State", "Private Endpoint", "Node Pools"}

	title := util.FormatColoredTitle(appCtx, fmt.Sprintf("Cluster: %s", c.DisplayName))
	p.PrintKeyValues(title, summary, order)
	fmt.Fprintln(appCtx.Stdout)

	if len(c.NodePools) > 0 {
		headers := []string{"Node Pool", "Version", "Shape", "Node Count"}
		rows := make([][]string, len(c.NodePools))
		for i, np := range c.NodePools {
			rows[i] = []string{
				np.DisplayName,
				np.KubernetesVersion,
				np.NodeShape,
				fmt.Sprintf("%d", np.NodeCount),
			}
		}
		tableTitle := util.FormatColoredTitle(appCtx, "Node Pools")
		p.PrintTable(tableTitle, headers, rows)
		fmt.Fprintln(appCtx.Stdout)
	}
}
