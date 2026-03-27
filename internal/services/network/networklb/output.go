package networklb

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	network "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintNetworkLoadBalancerInfo(nlb *network.NetworkLoadBalancer, appCtx *app.ApplicationContext, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(nlb)
	}

	title := util.FormatColoredTitle(appCtx, nlb.Name)

	if showAll {
		printAll(p, title, nlb)
	} else {
		printDefault(p, title, nlb)
	}

	return nil
}

func printDefault(p *printer.Printer, title string, nlb *network.NetworkLoadBalancer) {
	created := ""
	if nlb.Created != nil {
		created = nlb.Created.Format("2006-01-02")
	}
	vcn := "-"
	if nlb.VcnName != "" {
		vcn = nlb.VcnName
	} else if nlb.VcnID != "" {
		vcn = nlb.VcnID
	}
	data := map[string]string{
		"Name":           nlb.Name,
		"Created":        created,
		"IP Addresses":   strings.Join(nlb.IPAddresses, ", "),
		"State":          nlb.State,
		"Type":           nlb.Type,
		"VCN Name":       vcn,
		"Listeners":      formatListeners(nlb.Listeners),
		"Backend Health": formatBackendHealth(nlb.BackendHealth),
	}
	order := []string{"Name", "Created", "IP Addresses", "State", "Type", "VCN Name", "Listeners", "Backend Health"}
	p.PrintKeyValues(title, data, order)
}

func printAll(p *printer.Printer, title string, nlb *network.NetworkLoadBalancer) {
	created := ""
	if nlb.Created != nil {
		created = nlb.Created.Format("2006-01-02")
	}
	data := map[string]string{
		"Name":         nlb.Name,
		"Created":      created,
		"IP Addresses": strings.Join(nlb.IPAddresses, ", "),
		"State":        nlb.State,
		"OCID":         nlb.OCID,
		"Type":         nlb.Type,
		"VCN Name": func() string {
			if nlb.VcnName != "" {
				return nlb.VcnName
			}
			if nlb.VcnID != "" {
				return nlb.VcnID
			}
			return "-"
		}(),
		"Subnets":        strings.Join(nlb.Subnets, ", "),
		"NSGs":           strings.Join(nlb.NSGs, ", "),
		"Listeners":      formatListeners(nlb.Listeners),
		"Backend Health": formatBackendHealth(nlb.BackendHealth),
	}
	order := []string{"Name", "Created", "IP Addresses", "State", "OCID", "Type", "VCN Name", "Subnets", "NSGs", "Listeners", "Backend Health"}

	bsNames := make([]string, 0, len(nlb.BackendSets))
	for name := range nlb.BackendSets {
		bsNames = append(bsNames, name)
	}
	sort.Strings(bsNames)
	for i, name := range bsNames {
		bs := nlb.BackendSets[name]
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

func formatListeners(listeners map[string]string) string {
	var parts []string
	for name, backend := range listeners {
		parts = append(parts, fmt.Sprintf("%s → %s", name, backend))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}

func formatBackendHealth(health map[string]string) string {
	if len(health) == 0 {
		return ""
	}
	keys := make([]string, 0, len(health))
	for k := range health {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	const maxLines = 6
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

// PrintNetworkLoadBalancersInfo displays a list of network load balancers in a table or JSON with pagination support.
func PrintNetworkLoadBalancersInfo(nlbs []network.NetworkLoadBalancer, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, showAll bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse(p, nlbs, pagination)
	}

	if util.ValidateAndReportEmpty(nlbs, pagination, appCtx.Stdout) {
		return nil
	}

	for i := range nlbs {
		nlb := nlbs[i]
		if err := PrintNetworkLoadBalancerInfo(&nlb, appCtx, false, showAll); err != nil {
			return err
		}
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
