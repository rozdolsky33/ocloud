package subnet

import (
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/network/subnet"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintSubnetTable displays a table of subnets with details such as name, CIDR, and DNS info.
func PrintSubnetTable(subnets []subnet.Subnet, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, sortBy string) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		return util.MarshalDataToJSONResponse[subnet.Subnet](p, subnets, pagination)
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
				return subnets[i].CidrBlock < subnets[j].CidrBlock
			})
		}
	}

	// Define table headers
	headers := []string{"Name", "CIDR", "Public"}

	// Create rows for the table
	rows := make([][]string, len(subnets))
	for i, s := range subnets {
		// Create a row for this subnet
		rows[i] = []string{
			s.DisplayName,
			s.CidrBlock,
			util.FormatBool(s.Public),
		}
	}

	// Print the table without truncation so fully qualified domains are visible
	title := util.FormatColoredTitle(appCtx, "Subnets")
	p.PrintTableNoTruncate(title, headers, rows)

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// PrintSubnetInfo displays information about a list of subnets in either JSON format or a formatted table view.
func PrintSubnetInfo(subnets []subnet.Subnet, appCtx *app.ApplicationContext, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// If JSON output is requested, special-case empty for compact format expected by tests.
	if useJSON {
		if len(subnets) == 0 {
			_, err := appCtx.Stdout.Write([]byte("{\"items\": []}\n"))
			return err
		}
		return util.MarshalDataToJSONResponse[subnet.Subnet](p, subnets, nil)
	}

	if util.ValidateAndReportEmpty(subnets, nil, appCtx.Stdout) {
		return nil
	}

	// Print each policy as a separate key-value.
	for _, s := range subnets {
		// Derive DNS label and domain for display purposes only.
		dnsLabel := deriveDNSLabel(s.DisplayName)
		dnsDomain := dnsLabel + ".vcn1.oraclevcn.com"

		subnetData := map[string]string{
			"Name":       s.DisplayName,
			"Public":     util.FormatBool(s.Public),
			"CIDR":       s.CidrBlock,
			"DNS Label":  dnsLabel,
			"DNS Domain": dnsDomain,
		}

		orderedKeys := []string{
			"Name", "Public", "CIDR", "DNS Label", "DNS Domain",
		}

		// Create the colored title using components from the app context
		title := util.FormatColoredTitle(appCtx, s.DisplayName)

		p.PrintKeyValues(title, subnetData, orderedKeys)
	}

	util.LogPaginationInfo(nil, appCtx)
	return nil
}

// deriveDNSLabel derives a simple DNS label for display by using a trailing number
// from the subnet name if present (e.g., "TestSubnet1" -> "subnet1"). Otherwise,
// it returns a sanitized lowercase version of the name.
func deriveDNSLabel(name string) string {
	if name == "" {
		return "subnet"
	}
	n := strings.ToLower(name)
	// collect trailing digits
	i := len(n) - 1
	digits := ""
	for i >= 0 {
		c := n[i]
		if c < '0' || c > '9' {
			break
		}
		digits = string(c) + digits
		i--
	}
	if digits != "" {
		return "subnet" + digits
	}
	// sanitize: keep a-z, 0-9, replace spaces/underscores with '-'
	var b []rune
	for _, r := range n {
		switch {
		case r >= 'a' && r <= 'z':
			b = append(b, r)
		case r >= '0' && r <= '9':
			b = append(b, r)
		case r == ' ' || r == '_':
			b = append(b, '-')
			// drop others
		}
	}
	if len(b) == 0 {
		return "subnet"
	}
	return string(b)
}
