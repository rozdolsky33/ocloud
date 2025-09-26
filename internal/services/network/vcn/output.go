package vcn

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintVCNsInfo prints summaries for a slice of VCNs either as formatted text or as JSON.
// When useJSON is true it marshals the provided VCNs (respecting pagination if non-nil) to JSON;
// otherwise it prints a key/value summary for each VCN (OCID, state, CIDR blocks, IPv6 status,
// DNS label/domain, DHCP options, created date) and conditionally expands related sections
// (gateways, subnets, NSGs, route tables, security lists) according to the boolean flags.
// If pagination is provided it will be adjusted prior to output and pagination info is logged after printing.
func PrintVCNsInfo(vcns []domain.VCN, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[domain.VCN](p, vcns, pagination)
	}

	for _, v := range vcns {
		title := util.FormatColoredTitle(appCtx, v.DisplayName)
		cidrs := strings.Join(v.CidrBlocks, ", ")
		ipv6 := "Disabled"
		if v.Ipv6Enabled {
			ipv6 = "Enabled"
		}
		dhcp := strings.TrimSpace(v.DhcpOptions.DisplayName)
		if dhcp == "" {
			if strings.TrimSpace(v.DhcpOptionsID) != "" {
				dhcp = v.DhcpOptionsID
			} else {
				dhcp = "-"
			}
		} else if strings.TrimSpace(v.DhcpOptions.DomainNameType) != "" {
			dhcp = dhcp + " (" + v.DhcpOptions.DomainNameType + ")"
		}
		data := map[string]string{
			"OCID":               v.OCID,
			"State":              strings.ToUpper(v.LifecycleState),
			"CIDR Blocks":        cidrs,
			"IPv6":               ipv6,
			"DNS Label / Domain": strings.TrimSpace(strings.Join([]string{v.DnsLabel, v.DomainName}, " / ")),
			"DHCP Options":       dhcp,
			"Created":            v.TimeCreated.Format("2006-01-02"),
		}

		order := []string{"OCID", "State", "CIDR Blocks", "IPv6", "DNS Label / Domain", "DHCP Options", "Created"}
		p.PrintKeyValues(title, data, order)

		if gateways {
			printGateways(p, v.Gateways)
		}
		if subnets {
			printSubnets(p, v)
		}
		if nsgs {
			printNSGs(p, v.NSGs)
		}
		if routes {
			printRouteTables(p, v.RouteTables)
		}
		if securityLists {
			printSecurityLists(p, v.SecurityLists)
		}
	}
	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

//---------------------------------------------------------------------------------------------------------------------

// PrintVCNInfo prints the VCN's summary to the application's stdout or marshals it as JSON when useJSON is true.
// When rendered as text it displays core VCN fields (OCID, State, CIDR Blocks, IPv6, DNS Label / Domain, DHCP Options, Created)
// and, if enabled via flags, also prints Gateways, Subnets, Network Security Groups, Route Tables, and Security Lists.
// Returns any error encountered while marshaling to JSON.
func PrintVCNInfo(v domain.VCN, appCtx *app.ApplicationContext, useJSON, gateways, subnets, nsgs, routes, securityLists bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(v)
	}

	title := util.FormatColoredTitle(appCtx, v.DisplayName)
	cidrs := strings.Join(v.CidrBlocks, ", ")
	ipv6 := "Disabled"
	if v.Ipv6Enabled {
		ipv6 = "Enabled"
	}
	dhcp := strings.TrimSpace(v.DhcpOptions.DisplayName)
	if dhcp == "" {
		if strings.TrimSpace(v.DhcpOptionsID) != "" {
			dhcp = v.DhcpOptionsID
		} else {
			dhcp = "-"
		}
	} else if strings.TrimSpace(v.DhcpOptions.DomainNameType) != "" {
		dhcp = dhcp + " (" + v.DhcpOptions.DomainNameType + ")"
	}
	data := map[string]string{
		"OCID":               v.OCID,
		"State":              strings.ToUpper(v.LifecycleState),
		"CIDR Blocks":        cidrs,
		"IPv6":               ipv6,
		"DNS Label / Domain": strings.TrimSpace(strings.Join([]string{v.DnsLabel, v.DomainName}, " / ")),
		"DHCP Options":       dhcp,
		"Created":            v.TimeCreated.Format("2006-01-02"),
	}

	order := []string{"OCID", "State", "CIDR Blocks", "IPv6", "DNS Label / Domain", "DHCP Options", "Created"}
	p.PrintKeyValues(title, data, order)

	if gateways {
		printGateways(p, v.Gateways)
	}
	if subnets {
		printSubnets(p, v)
	}
	if nsgs {
		printNSGs(p, v.NSGs)
	}
	if routes {
		printRouteTables(p, v.RouteTables)
	}
	if securityLists {
		printSecurityLists(p, v.SecurityLists)
	}

	return nil
}

// printGateways prints a "Gateways" table with columns "Type" and "Details" for the provided gateways.
// If the slice is empty, the function returns without producing output.

func printGateways(p *printer.Printer, gateways []domain.Gateway) {
	if len(gateways) == 0 {
		return
	}
	p.PrintTable("Gateways", []string{"Type", "Details"}, toGatewayRows(gateways))
}

// printSubnets prints a table of the VCN's subnets when the VCN has any subnets.
// The table is titled "Subnets" and includes columns: Name, CIDR, Publicity, Route Table, and SecLists (no truncation).
func printSubnets(p *printer.Printer, v domain.VCN) {
	subnets := v.Subnets
	if len(subnets) == 0 {
		return
	}
	headers := []string{"Name", "CIDR", "Publicity", "Route Table", "SecLists"}
	p.PrintTableNoTruncate("Subnets", headers, toSubnetRows(v))
}

// The table contains the columns "Name" and "State".
func printNSGs(p *printer.Printer, nsgs []domain.NSG) {
	if len(nsgs) == 0 {
		return
	}
	headers := []string{"Name", "State"}
	p.PrintTableNoTruncate("Network Security Groups", headers, toNSGRows(nsgs))
}

// printRouteTables prints a "Route Tables" table showing each route table's name and lifecycle state when the provided slice is non-empty.
func printRouteTables(p *printer.Printer, rts []domain.RouteTable) {
	if len(rts) == 0 {
		return
	}
	headers := []string{"Name", "State"}
	p.PrintTableNoTruncate("Route Tables", headers, toRouteTableRows(rts))
}

// printSecurityLists prints a "Security Lists" table with columns "Name" and "State" using the provided printer if the slice is non-empty.
func printSecurityLists(p *printer.Printer, sls []domain.SecurityList) {
	if len(sls) == 0 {
		return
	}
	headers := []string{"Name", "State"}
	p.PrintTableNoTruncate("Security Lists", headers, toSecurityListRows(sls))
}

// toGatewayRows groups gateways by type and returns rows suitable for tabular display.
// It aggregates gateway display names by type and returns a slice of rows where each
// row is [label, comma-separated names]. Gateway types with no entries are omitted.
func toGatewayRows(gateways []domain.Gateway) [][]string {
	var (
		internet []string
		nat      []string
		service  []string
		drg      []string
		lpg      []string
	)
	for _, gw := range gateways {
		switch gw.Type {
		case "Internet":
			internet = append(internet, gw.DisplayName)
		case "NAT":
			nat = append(nat, gw.DisplayName)
		case "Service":
			service = append(service, gw.DisplayName)
		case "DRG":
			drg = append(drg, gw.DisplayName)
		case "Local Peering":
			lpg = append(lpg, gw.DisplayName)
		}
	}
	var rows [][]string
	if len(internet) > 0 {
		rows = append(rows, []string{"Internet", strings.Join(internet, ", ")})
	}
	if len(nat) > 0 {
		rows = append(rows, []string{"NAT", strings.Join(nat, ", ")})
	}
	if len(service) > 0 {
		rows = append(rows, []string{"Service GW", strings.Join(service, ", ")})
	}
	if len(drg) > 0 {
		rows = append(rows, []string{"DRG", strings.Join(drg, ", ")})
	}
	if len(lpg) > 0 {
		rows = append(rows, []string{"LPG Peers", strings.Join(lpg, ", ")})
	}
	return rows
}

// toSubnetRows builds table rows for each subnet in the given VCN.
// 
// Each returned row contains these columns in order: subnet display name, CIDR block,
// publicity ("PUBLIC" or "PRIVATE"), resolved route table name (or OCID) with line-wrapping,
// and resolved security list names (or OCIDs) with line-wrapping.
func toSubnetRows(v domain.VCN) [][]string {
	subnets := v.Subnets
	rows := make([][]string, len(subnets))
	for i, s := range subnets {
		rt := lookupRouteTableName(v, s.RouteTableID)
		rt = strings.Join(util.SplitTextByMaxWidth(rt), "\n")
		sl := lookupSecurityListNames(v, s.SecurityListIDs)
		sl = strings.Join(util.SplitTextByMaxWidth(sl), "\n")
		rows[i] = []string{
			s.DisplayName,
			s.CidrBlock,
			formatPublicity(s.Public),
			rt,
			sl,
		}
	}
	return rows
}

// formatPublicity returns "PUBLIC" when public is true and "PRIVATE" when public is false.

func formatPublicity(public bool) string {
	if public {
		return "PUBLIC"
	}
	return "PRIVATE"
}

// lookupRouteTableName resolves a route table's display name from its OCID within the given VCN.
// If `id` is empty or whitespace it returns "-".
// If a route table with a matching OCID is found and has a non-empty display name, that display name is returned;
// otherwise the provided `id` is returned.
func lookupRouteTableName(v domain.VCN, id string) string {
	if strings.TrimSpace(id) == "" {
		return "-"
	}
	for _, rt := range v.RouteTables {
		if rt.OCID == id {
			if strings.TrimSpace(rt.DisplayName) != "" {
				return rt.DisplayName
			}
			return id
		}
	}
	return id
}

// lookupSecurityListNames builds a comma-separated string of display names for the given security list OCIDs.
// If a security list has an empty display name, its OCID is used in its place. If the ids slice is empty, it returns "-".
func lookupSecurityListNames(v domain.VCN, ids []string) string {
	if len(ids) == 0 {
		return "-"
	}
	nameByID := make(map[string]string, len(v.SecurityLists))
	for _, sl := range v.SecurityLists {
		nameByID[sl.OCID] = sl.DisplayName
	}
	var names []string
	for _, id := range ids {
		name := nameByID[id]
		if strings.TrimSpace(name) == "" {
			name = id
		}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// Each row contains the NSG display name and the lifecycle state in uppercase.
func toNSGRows(nsgs []domain.NSG) [][]string {
	rows := make([][]string, len(nsgs))
	for i, n := range nsgs {
		rows[i] = []string{n.DisplayName, strings.ToUpper(n.LifecycleState)}
	}
	return rows
}

// toRouteTableRows builds table rows for the provided route tables.
// Each row contains the route table's display name and its lifecycle state in uppercase, returned as a slice of string slices.
func toRouteTableRows(rts []domain.RouteTable) [][]string {
	rows := make([][]string, len(rts))
	for i, r := range rts {
		rows[i] = []string{r.DisplayName, strings.ToUpper(r.LifecycleState)}
	}
	return rows
}

// toSecurityListRows converts a slice of SecurityList into table rows suitable for display.
// Each returned row contains the security list's DisplayName and its LifecycleState in uppercase.
func toSecurityListRows(sls []domain.SecurityList) [][]string {
	rows := make([][]string, len(sls))
	for i, s := range sls {
		rows[i] = []string{s.DisplayName, strings.ToUpper(s.LifecycleState)}
	}
	return rows
}
