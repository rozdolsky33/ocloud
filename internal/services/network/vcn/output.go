package vcn

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintVCNSummary prints the VCN summary view or JSON if requested.
func PrintVCNSummary(v *VCNDTO, appCtx *app.ApplicationContext, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(v)
	}

	// Build title: <tenancy>: <compartment>: <vcn-name> (<region if known>)
	title := util.FormatColoredTitle(appCtx, v.DisplayName)

	// Prepare key-value data
	cidrs := strings.Join(v.CidrBlocks, ", ")
	ipv6 := "Disabled"
	if v.Ipv6Enabled {
		ipv6 = "Enabled"
	}

	// OCID short form: prefix + ellipsis + last 4
	ocidShort := shortenOCID(v.OCID)
	ocidDisplay := fmt.Sprintf("%s (%s)", ocidShort, shortIDToken(v.OCID))

	compVal := v.CompartmentID
	if v.CompartmentName != "" {
		compVal = fmt.Sprintf("%s (%s)", v.CompartmentName, shortenOCID(v.CompartmentID))
	} else if appCtx.CompartmentID == v.CompartmentID && appCtx.CompartmentName != "" {
		compVal = fmt.Sprintf("%s (%s)", appCtx.CompartmentName, shortenOCID(v.CompartmentID))
	} else {
		compVal = fmt.Sprintf("%s", shortenOCID(v.CompartmentID))
	}

	dhcpVal := "-"
	if v.DhcpOptionsName != "" || v.DhcpCustomDNS != "" {
		label := v.DhcpOptionsName
		if label == "" {
			label = v.DhcpOptionsID
		}
		labelShort := label
		if label == v.DhcpOptionsID {
			labelShort = shortenOCID(v.DhcpOptionsID)
		}
		details := ""
		if v.DhcpCustomDNS != "" {
			details = fmt.Sprintf(" (custom DNS: %s)", v.DhcpCustomDNS)
		}
		dhcpVal = fmt.Sprintf("%s%s", labelShort, details)
	}

	tags := flattenTags(v.FreeformTags)
	if tags == "" {
		tags = "-"
	}

	data := map[string]string{
		"Name":               v.DisplayName,
		"OCID":               ocidDisplay,
		"State":              strings.ToUpper(v.LifecycleState),
		"Compartment":        compVal,
		"CIDR Blocks":        cidrs,
		"IPv6":               ipv6,
		"DNS Label / Domain": strings.TrimSpace(strings.Join([]string{v.DnsLabel, v.DomainName}, " / ")),
		"DHCP Options":       dhcpVal,
		"Created":            v.TimeCreated,
		"Tags":               tags,
	}

	order := []string{"Name", "OCID", "State", "Compartment", "CIDR Blocks", "IPv6", "DNS Label / Domain", "DHCP Options", "Created", "Tags"}
	p.PrintKeyValues(title, data, order)
	util.LogPaginationInfo(nil, appCtx)
	return nil
}

func PrintVCNsTable(vcns []*VCNDTO, appCtx *app.ApplicationContext, useJSON bool) error {
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
			shortenOCID(v.OCID),
			strings.ToUpper(v.LifecycleState),
			strings.Join(v.CidrBlocks, ", "),
			flattenTags(v.FreeformTags),
		}
		data = append(data, row)
	}

	p.PrintTable(title, headers, data)
	util.LogPaginationInfo(nil, appCtx)

	return nil
}

func shortenOCID(ocid string) string {
	if ocid == "" {
		return ""
	}
	if len(ocid) <= 16 {
		return ocid
	}
	return ocid[:12] + "…" + ocid[len(ocid)-4:]
}

// shortIDToken returns resource code like vcn-xxxx from the OCID if possible
func shortIDToken(ocid string) string {
	// OCI OCID often contains a friendly token at the end after a dot, but not reliable.
	// We'll synthesize vcn-<last4> when ocid contains ".vcn" pattern.
	last4 := ""
	if len(ocid) >= 4 {
		last4 = ocid[len(ocid)-4:]
	}
	kind := "res"
	if strings.Contains(ocid, ".vcn") || strings.Contains(ocid, "ocid1.vcn") {
		kind = "vcn"
	}
	return fmt.Sprintf("%s-%s", kind, last4)
}

func flattenTags(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}
	parts := make([]string, 0, len(tags))
	for k, v := range tags {
		// color value for a nicer look
		colored := text.Colors{text.FgYellow}.Sprintf("%s=%s", k, v)
		parts = append(parts, colored)
	}
	// stable order is not critical; join with comma-space
	return strings.Join(parts, ", ")
}
