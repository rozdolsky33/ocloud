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
	rp := strings.Join(lb.RoutingPolicies, ", ")
	if rp == "" {
		rp = "-"
	}
	useSSL := "No"
	if lb.UseSSL {
		useSSL = "Yes"
	}
	data := map[string]string{
		"Name":           lb.Name,
		"Shape":          lb.Shape,
		"Created":        created,
		"IP Addresses":   strings.Join(lb.IPAddresses, ", "),
		"State":          lb.State,
		"Listeners":      formatListeners(lb.Listeners, false),
		"Backend Health": formatBackendHealth(lb.BackendHealth),
		"Routing Policy": rp,
		"Use SSL":        useSSL,
	}
	order := []string{"Name", "Shape", "Created", "IP Addresses", "State", "Listeners", "Backend Health", "Routing Policy", "Use SSL"}
	p.PrintKeyValues(title, data, order)
}

func printAll(p *printer.Printer, title string, lb *network.LoadBalancer) {
	created := ""
	if lb.Created != nil {
		created = lb.Created.String()
	}
	data := map[string]string{
		"Name":           lb.Name,
		"Shape":          lb.Shape,
		"Created":        created,
		"IP Addresses":   strings.Join(lb.IPAddresses, ", "),
		"State":          lb.State,
		"OCID":           lb.OCID,
		"Type":           lb.Type,
		"Subnets":        strings.Join(lb.Subnets, ", "),
		"NSGs":           strings.Join(lb.NSGs, ", "),
		"Listeners":      formatListeners(lb.Listeners, true),
		"Backend Health": formatBackendHealth(lb.BackendHealth),
		"Routing Policy": func() string {
			if len(lb.RoutingPolicies) == 0 {
				return "-"
			}
			return strings.Join(lb.RoutingPolicies, ", ")
		}(),
		"Use SSL": func() string {
			if lb.UseSSL {
				return "Yes"
			}
			return "No"
		}(),
		"SSL Certificates": formatCertificates(lb.SSLCertificates),
	}
	order := []string{"Name", "Shape", "Created", "IP Addresses", "State", "OCID", "Type", "Subnets", "NSGs", "Listeners", "Backend Health", "Routing Policy", "Use SSL", "SSL Certificates"}

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

	p.PrintKeyValuesNoTruncate(title, data, order)
}

func formatListeners(listeners map[string]string, includeNames bool) string {
	var parts []string
	if includeNames {
		for name, backend := range listeners {
			normalized := normalizeListenerValue(backend)
			parts = append(parts, fmt.Sprintf("%s → %s", name, normalized))
		}
	} else {
		// Default view: omit listener names, show only protocol:port → backendset
		for _, backend := range listeners {
			normalized := normalizeListenerValue(backend)
			parts = append(parts, normalized)
		}
	}
	return strings.Join(parts, "\n")
}

// normalizeListenerValue ensures we display correct scheme labels for common ports
// and previously mislabeled values coming from upstream mapping. It operates on
// strings like "http:8443 → backendset" and fixes them to "https:8443 → backendset".
func normalizeListenerValue(s string) string {
	// Split on the arrow to isolate the left side (proto:port)
	leftRight := strings.SplitN(s, " → ", 2)
	left := leftRight[0]
	right := ""
	if len(leftRight) == 2 {
		right = leftRight[1]
	}
	// Expect left to be proto:port
	lp := strings.SplitN(left, ":", 2)
	if len(lp) != 2 {
		return s
	}
	proto := strings.ToLower(strings.TrimSpace(lp[0]))
	portStr := strings.TrimSpace(lp[1])
	// Extract numeric port (strip anything after space just in case)
	if i := strings.IndexByte(portStr, ' '); i >= 0 {
		portStr = portStr[:i]
	}
	// Force schemes for common ports
	switch portStr {
	case "443", "8443":
		proto = "https"
	case "80":
		proto = "http"
	}
	leftFixed := fmt.Sprintf("%s:%s", proto, portStr)
	if right != "" {
		return leftFixed + " → " + right
	}
	return leftFixed
}

func formatBackendHealth(health map[string]string) string {
	if len(health) == 0 {
		return ""
	}
	// Sort backend set names for deterministic output, then limit lines to keep the table compact
	keys := make([]string, 0, len(health))
	for k := range health {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	const maxLines = 6 // show up to 6 entries, then summarize the rest
	var parts []string
	limit := len(keys)
	if limit > maxLines {
		limit = maxLines
	}
	for i := 0; i < limit; i++ {
		k := keys[i]
		parts = append(parts, fmt.Sprintf("%s: %s", k, health[k]))
	}
	if len(keys) > maxLines {
		parts = append(parts, fmt.Sprintf("… (+%d more)", len(keys)-maxLines))
	}
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

func formatCertificates(certs []string) string {
	if len(certs) == 0 {
		return ""
	}
	return strings.Join(certs, "\n")
}
