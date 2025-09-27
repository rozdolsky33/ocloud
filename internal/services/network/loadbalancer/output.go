package loadbalancer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	network "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintLoadBalancerInfo(lb *network.LoadBalancer, appCtx *app.ApplicationContext, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(lb)
	}

	title := util.FormatColoredTitle(appCtx, lb.Name)

	if showAll {
		printAll(p, title, lb)
	} else {
		printDefault(p, title, lb)
	}

	return nil
}

func printDefault(p *printer.Printer, title string, lb *network.LoadBalancer) {
	created := ""
	if lb.Created != nil {
		created = lb.Created.String()
	}
	data := map[string]string{
		"Name":           lb.Name,
		"Shape":          lb.Shape,
		"Created":        created,
		"IP":             strings.Join(lb.IPAddresses, ", "),
		"State":          lb.State,
		"Listeners":      formatListeners(lb.Listeners, false),
		"Backend Health": formatBackendHealth(lb.BackendHealth),
	}
	order := []string{"Name", "Shape", "Created", "IP Addresses", "State", "Listeners", "Backend Health"}
	p.PrintKeyValues(title, data, order)
}

func printAll(p *printer.Printer, title string, lb *network.LoadBalancer) {
	created := ""
	if lb.Created != nil {
		created = lb.Created.String()
	}
	data := map[string]string{
		"Name":             lb.Name,
		"Shape":            lb.Shape,
		"Created":          created,
		"IP":               strings.Join(lb.IPAddresses, ", "),
		"State":            lb.State,
		"OCID":             lb.OCID,
		"Type":             lb.Type,
		"Subnets":          strings.Join(lb.Subnets, ", "),
		"NSGs":             strings.Join(lb.NSGs, ", "),
		"Listeners":        formatListeners(lb.Listeners, true),
		"Backend Health":   formatBackendHealth(lb.BackendHealth),
		"SSL Certificates": strings.Join(lb.SSLCertificates, ", "),
	}
	order := []string{"Name", "Shape", "Created", "IP Addresses", "State", "OCID", "Type", "Subnets", "NSGs", "Listeners", "Backend Health", "SSL Certificates"}

	// Include backend set summaries as additional key-value entries (no separate tables)
	// To avoid truncating long backend set names in the Key column, we print a short key
	// like "Backend Set 1" and include the full backend set name as the first line of the value.
	// Also, sort backend set names for a stable output order.
	bsNames := make([]string, 0, len(lb.BackendSets))
	for name := range lb.BackendSets {
		bsNames = append(bsNames, name)
	}
	sort.Strings(bsNames)
	for i, name := range bsNames {
		bs := lb.BackendSets[name]
		key := fmt.Sprintf("Backend Set %d", i+1)
		val := fmt.Sprintf("%s\nPolicy: %s, HC: %s", name, bs.Policy, bs.Health)
		if len(bs.Backends) > 0 {
			parts := make([]string, 0, len(bs.Backends))
			for _, b := range bs.Backends {
				parts = append(parts, fmt.Sprintf("%s:%d (%s)", b.Name, b.Port, b.Status))
			}
			val = val + "\nBackends: " + strings.Join(parts, ", ")
		}
		data[key] = val
		order = append(order, key)
	}

	p.PrintKeyValues(title, data, order)
}

func formatListeners(listeners map[string]string, includeNames bool) string {
	var parts []string
	if includeNames {
		for name, backend := range listeners {
			parts = append(parts, fmt.Sprintf("%s → %s", name, backend))
		}
	} else {
		// Default view: omit listener names, show only protocol:port → backendset
		for _, backend := range listeners {
			parts = append(parts, backend)
		}
	}
	return strings.Join(parts, "\n")
}

func formatBackendHealth(health map[string]string) string {
	var parts []string
	for backend, status := range health {
		parts = append(parts, fmt.Sprintf("%s: %s", backend, status))
	}
	// Use newline to match Listeners formatting and avoid overly wide tables
	return strings.Join(parts, "\n")
}

// PrintLoadBalancersInfo displays a list of load balancers in a table or JSON with pagination support.
func PrintLoadBalancersInfo(lbs []network.LoadBalancer, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse(p, lbs, pagination)
	}

	if util.ValidateAndReportEmpty(lbs, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each load balancer as a key-value block similar to instance output
	for i := range lbs {
		lb := lbs[i]
		// Reuse single-resource printer to ensure consistent formatting
		if err := PrintLoadBalancerInfo(&lb, appCtx, false, showAll); err != nil {
			return err
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
