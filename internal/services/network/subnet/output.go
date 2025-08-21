package subnet

import (
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintSubnetTable displays a table of subnets with details such as name, CIDR, and DNS info.
func PrintSubnetTable(subnets []Subnet, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, sortBy string) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[Subnet](p, subnets, pagination)
	}

	if util.ValidateAndReportEmpty(subnets, pagination, appCtx.Stdout) {
		return nil
	}

	// Sort subnets based on sortBy parameter
	if sortBy != "" {
		sortBy = strings.ToLower(sortBy)
		switch sortBy {
		case "name":
			sort.Slice(subnets, func(i, j int) bool {
				return strings.ToLower(subnets[i].DisplayName) < strings.ToLower(subnets[j].DisplayName)
			})
		case "cidr":
			sort.Slice(subnets, func(i, j int) bool {
				return subnets[i].CIDRBlock < subnets[j].CIDRBlock
			})
		}
	}

	// Define table headers
	headers := []string{"Name", "CIDR", "Public IP", "DNS Label", "Subnet Domain"}

	// Create rows for the table
	rows := make([][]string, len(subnets))
	for i, subnet := range subnets {
		// Determine if public IP is allowed
		publicIPAllowed := "No"
		if !subnet.ProhibitPublicIPOnVnic {
			publicIPAllowed = "Yes"
		}

		// Create a row for this subnet
		rows[i] = []string{
			subnet.DisplayName,
			subnet.CIDRBlock,
			publicIPAllowed,
			subnet.DNSLabel,
			subnet.SubnetDomainName,
		}
	}

	// Print the table without truncation so fully qualified domains are visible
	title := util.FormatColoredTitle(appCtx, "Subnets")
	p.PrintTableNoTruncate(title, headers, rows)

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// PrintSubnetInfo displays information about a list of subnets in either JSON format or a formatted table view.
func PrintSubnetInfo(subnets []Subnet, appCtx *app.ApplicationContext, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// If JSON output is requested, special-case empty for compact format expected by tests.
	if useJSON {
		if len(subnets) == 0 {
			// Write compact JSON: {"items": []}
			_, err := appCtx.Stdout.Write([]byte("{\"items\": []}\n"))
			return err
		}
		return util.MarshalDataToJSONResponse[Subnet](p, subnets, nil)
	}

	if util.ValidateAndReportEmpty(subnets, nil, appCtx.Stdout) {
		return nil
	}

	// Print each policy as a separate key-value.
	for _, subnet := range subnets {
		publicIPAllowed := "No"
		if !subnet.ProhibitPublicIPOnVnic {
			publicIPAllowed = "Yes"
		}
		subnetData := map[string]string{
			"Name":          subnet.DisplayName,
			"Public IP":     publicIPAllowed,
			"CIDR":          subnet.CIDRBlock,
			"DNS Label":     subnet.DNSLabel,
			"Subnet Domain": subnet.SubnetDomainName,
		}

		orderedKeys := []string{
			"Name", "Public IP", "CIDR", "DNS Label", "Subnet Domain",
		}

		// Create the colored title using components from the app context
		title := util.FormatColoredTitle(appCtx, subnet.DisplayName)

		p.PrintKeyValues(title, subnetData, orderedKeys)
	}

	util.LogPaginationInfo(nil, appCtx)
	return nil
}
