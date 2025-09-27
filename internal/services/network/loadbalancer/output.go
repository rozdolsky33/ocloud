package loadbalancer

import (
	"fmt"
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

	title := fmt.Sprintf("%s: %s: %s", appCtx.TenancyName, appCtx.CompartmentName, lb.Name)

	if showAll {
		printAll(p, title, lb)
	} else {
		printDefault(p, title, lb)
	}

	return nil
}

func printDefault(p *printer.Printer, title string, lb *network.LoadBalancer) {
	data := map[string]string{
		"Name":           lb.Name,
		"State":          lb.State,
		"Type":           lb.Type,
		"IP Addresses":   strings.Join(lb.IPAddresses, ", "),
		"Shape":          lb.Shape,
		"Listeners":      formatListeners(lb.Listeners),
		"Backend Health": formatBackendHealth(lb.BackendHealth),
	}
	order := []string{"Name", "State", "Type", "IP Addresses", "Shape", "Listeners", "Backend Health"}
	p.PrintKeyValues(title, data, order)
}

func printAll(p *printer.Printer, title string, lb *network.LoadBalancer) {
	created := ""
	if lb.Created != nil {
		created = lb.Created.Format("2006-01-02")
	}
	data := map[string]string{
		"Name":         lb.Name,
		"OCID":         lb.OCID,
		"State":        lb.State,
		"Type":         lb.Type,
		"Shape":        lb.Shape,
		"IP Addresses": strings.Join(lb.IPAddresses, ", "),
		"Subnets":      strings.Join(lb.Subnets, ", "),
		"NSGs":         strings.Join(lb.NSGs, ", "),
		"Created":      created,
		"Listeners":    formatListeners(lb.Listeners),
	}
	order := []string{"Name", "OCID", "State", "Type", "Shape", "IP Addresses", "Subnets", "NSGs", "Created", "Listeners"}
	p.PrintKeyValues(title, data, order)

	// Backend Sets
	if len(lb.BackendSets) > 0 {
		for name, bs := range lb.BackendSets {
			bsTitle := fmt.Sprintf("Backend Set: %s", name)
			rows := [][]string{}
			// the header row will be provided to PrintTable
			for _, b := range bs.Backends {
				rows = append(rows, []string{b.Name, fmt.Sprintf("%d", b.Port), b.Status})
			}
			// Print policy and health as a preface rowless table by including them in the title
			preface := fmt.Sprintf("Policy: %s, HC: %s", bs.Policy, bs.Health)
			p.PrintTable(bsTitle+"\n"+preface, []string{"Backend", "Port", "Status"}, rows)
		}
	}

	// SSL Certificates
	if len(lb.SSLCertificates) > 0 {
		rows := make([][]string, 0, len(lb.SSLCertificates))
		for _, c := range lb.SSLCertificates {
			// We expect pre-formatted strings like "name (Expires: YYYY-MM-DD)" in the domain layer
			rows = append(rows, []string{c})
		}
		p.PrintTable("SSL Certificates", []string{"Certificate"}, rows)
	}
}

func formatListeners(listeners map[string]string) string {
	var parts []string
	for name, backend := range listeners {
		parts = append(parts, fmt.Sprintf("%s → %s", name, backend))
	}
	return strings.Join(parts, "\n")
}

func formatBackendHealth(health map[string]string) string {
	var parts []string
	for backend, status := range health {
		parts = append(parts, fmt.Sprintf("%s: %s", backend, status))
	}
	return strings.Join(parts, ", ")
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

	headers := []string{"Name", "State", "Type", "IPs", "Shape", "Created"}
	if showAll {
		headers = append(headers, "Subnets", "NSGs")
	}

	rows := make([][]string, 0, len(lbs))
	for _, lb := range lbs {
		created := ""
		if lb.Created != nil {
			created = lb.Created.Format("2006-01-02")
		}
		ips := strings.Join(lb.IPAddresses, ", ")
		row := []string{lb.Name, lb.State, lb.Type, ips, lb.Shape, created}
		if showAll {
			row = append(row, strings.Join(lb.Subnets, ", "), strings.Join(lb.NSGs, ", "))
		}
		rows = append(rows, row)
	}

	title := util.FormatColoredTitle(appCtx, "Load Balancers")
	p.PrintTable(title, headers, rows)

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
