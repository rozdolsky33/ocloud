package vcn

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintVCNsInfo prints the VCN summary view or JSON if requested.
func PrintVCNsInfo(vcns []vcn.VCN, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON, gateways, subnets bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[vcn.VCN](p, vcns, pagination)
	}

	for _, v := range vcns {
		title := util.FormatColoredTitle(appCtx, v.DisplayName)
		cidrs := strings.Join(v.CidrBlocks, ", ")
		ipv6 := "Disabled"
		if v.Ipv6Enabled {
			ipv6 = "Enabled"
		}
		// Determine DHCP Options label with fallback to ID and optional domain type
		dhcp := strings.TrimSpace(v.DhcpOptions.DisplayName)
		if dhcp == "" {
			if strings.TrimSpace(v.DhcpOptionsID) != "" {
				dhcp = v.DhcpOptionsID
			} else {
				dhcp = "-"
			}
		} else if strings.TrimSpace(v.DhcpOptions.DomainNameType) != "" {
			// Show the domain name type if we have it
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
	}
	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func printGateways(p *printer.Printer, gateways []vcn.Gateway) {
	if len(gateways) == 0 {
		return
	}
	p.PrintTable("Gateways", []string{"Type", "Details"}, toGatewayRows(gateways))
}

func printSubnets(p *printer.Printer, v vcn.VCN) {
	subnets := v.Subnets
	if len(subnets) == 0 {
		return
	}
	headers := []string{"Name", "CIDR", "Publicity", "Route Table", "SecLists", "NSGs", "Egress Path"}
	// Use non-truncating table to ensure full information is visible
	p.PrintTableNoTruncate("Subnets", headers, toSubnetRows(v))
}

func toGatewayRows(gateways []vcn.Gateway) [][]string {
	// Group gateways by type and format details without OCIDs
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
			internet = append(internet, gw.DisplayName+" (present)")
		case "NAT":
			nat = append(nat, gw.DisplayName+" (present)")
		case "Service":
			service = append(service, gw.DisplayName+" (present)")
		case "DRG":
			drg = append(drg, gw.DisplayName+" (attached)")
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

func toSubnetRows(v vcn.VCN) [][]string {
	subnets := v.Subnets
	rows := make([][]string, len(subnets))
	for i, s := range subnets {
		rt := lookupRouteTableName(v, s.RouteTableID)
		rt = strings.Join(util.SplitTextByMaxWidth(rt), "\n")
		sl := lookupSecurityListNames(v, s.SecurityListIDs)
		sl = strings.Join(util.SplitTextByMaxWidth(sl), "\n")
		nsg := lookupNSGNames(v, s.NSGIDs)
		nsg = strings.Join(util.SplitTextByMaxWidth(nsg), "\n")
		rows[i] = []string{
			s.DisplayName,
			s.CidrBlock,
			formatPublicity(s.Public),
			rt,
			sl,
			nsg,
			estimateEgressPath(v, s.RouteTableID),
		}
	}
	return rows
}

// --- helpers ---

func formatPublicity(public bool) string {
	if public {
		return "PUBLIC"
	}
	return "PRIVATE"
}

func lookupRouteTableName(v vcn.VCN, id string) string {
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

func lookupSecurityListNames(v vcn.VCN, ids []string) string {
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

func lookupNSGNames(v vcn.VCN, ids []string) string {
	// If subnet doesn't explicitly list NSG IDs, fall back to showing all NSGs present in the VCN
	if len(ids) == 0 {
		if len(v.NSGs) == 0 {
			return "-"
		}
		var all []string
		for _, n := range v.NSGs {
			label := strings.TrimSpace(n.DisplayName)
			if label == "" {
				label = n.OCID
			}
			all = append(all, label)
		}
		return strings.Join(all, ", ")
	}
	nameByID := make(map[string]string, len(v.NSGs))
	for _, n := range v.NSGs {
		nameByID[n.OCID] = n.DisplayName
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

func estimateEgressPath(v vcn.VCN, routeTableID string) string {
	name := strings.ToLower(lookupRouteTableName(v, routeTableID))
	if strings.Contains(name, "igw") || strings.Contains(name, "public") {
		return "IGW"
	}
	if strings.Contains(name, "nat") || strings.Contains(name, "private") {
		return "NAT"
	}
	if strings.Contains(name, "drg") || strings.Contains(name, "hub") || strings.Contains(name, "db") {
		return "DRG"
	}
	return "-"
}
