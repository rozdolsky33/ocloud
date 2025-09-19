package vcn

import (
	"fmt"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintVCNSummary prints the VCN summary view or JSON if requested.
func PrintVCNSummary(v *VCN, appCtx *app.ApplicationContext, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(v)
	}

	// Build title: <tenancy>: <compartment>: <vcn-name> (<region if known>)
	title := util.FormatColoredTitle(appCtx, v.DisplayName)

	cidrs := strings.Join(v.CidrBlocks, ", ")
	ipv6 := "Disabled"
	if v.Ipv6Enabled {
		ipv6 = "Enabled"
	}

	dhcp := v.DhcpOptionsID
	if v.DhcpOptions.DisplayName != "" {
		dhcp = fmt.Sprintf("%s (%s)", v.DhcpOptions.DisplayName, v.DhcpOptions.OCID)
	}

	data := map[string]string{
		"Name":               v.DisplayName,
		"OCID":               v.OCID,
		"State":              strings.ToUpper(v.LifecycleState),
		"Compartment":        appCtx.CompartmentName,
		"CIDR Blocks":        cidrs,
		"IPv6":               ipv6,
		"DNS Label / Domain": strings.TrimSpace(strings.Join([]string{v.DnsLabel, v.DomainName}, " / ")),
		"DHCP Options":       dhcp,
		"Created":            v.TimeCreated.Format("2006-01-02"),
	}

	order := []string{"Name", "OCID", "State", "Compartment", "CIDR Blocks", "IPv6", "DNS Label / Domain", "DHCP Options", "Created"}
	p.PrintKeyValues(title, data, order)
	util.LogPaginationInfo(nil, appCtx)
	return nil
}

// PrintVCNsInfo prints the VCN summary view or JSON if requested.
func PrintVCNsInfo(vcns []VCN, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[VCN](p, vcns, pagination)
	}

	for _, v := range vcns {
		title := util.FormatColoredTitle(appCtx, v.DisplayName)
		cidrs := strings.Join(v.CidrBlocks, ", ")
		ipv6 := "Disabled"
		if v.Ipv6Enabled {
			ipv6 = "Enabled"
		}
		data := map[string]string{
			"OCID":               v.OCID,
			"State":              strings.ToUpper(v.LifecycleState),
			"CIDR Blocks":        cidrs,
			"IPv6":               ipv6,
			"DNS Label / Domain": strings.TrimSpace(strings.Join([]string{v.DnsLabel, v.DomainName}, " / ")),
			"DHCP Options":       v.DhcpOptions.DisplayName,
			"Created":            v.TimeCreated.Format("2006-01-02"),
		}

		order := []string{"OCID", "State", "CIDR Blocks", "IPv6", "DNS Label / Domain", "DHCP Options", "Created"}
		p.PrintKeyValues(title, data, order)
	}
	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

func PrintVCNsTable(vcns []*VCN, appCtx *app.ApplicationContext, useJSON bool) error {
	p := printer.New(appCtx.Stdout)
	if useJSON {
		return p.MarshalToJSON(vcns)
	}

	title := util.FormatColoredTitle(appCtx, "VCNs")
	headers := []string{"Name", "OCID", "State", "CIDR Blocks", "Tags"}

	data := make([][]string, 0, len(vcns))
	for _, v := range vcns {
		row := []string{
			v.DisplayName,
			v.OCID,
			strings.ToUpper(v.LifecycleState),
			strings.Join(v.CidrBlocks, ", "),
		}
		data = append(data, row)
	}

	p.PrintTable(title, headers, data)
	util.LogPaginationInfo(nil, appCtx)

	return nil
}
