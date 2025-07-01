package subnet

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"sort"
	"strings"
)

func PrintSubnetInfo(subnets []Subnet, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, sortBy string) error {

	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		// Special case for empty compartments list - return an empty object
		if len(subnets) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
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
				return strings.ToLower(subnets[i].Name) < strings.ToLower(subnets[j].Name)
			})
		case "cidr":
			sort.Slice(subnets, func(i, j int) bool {
				return subnets[i].CIDR < subnets[j].CIDR
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
			subnet.Name,
			subnet.CIDR,
			publicIPAllowed,
			subnet.DNSLabel,
			subnet.SubnetDomainName,
		}
	}

	// Print the table
	title := util.FormatColoredTitle(appCtx, "Subnets")
	p.PrintTable(title, headers, rows)

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
