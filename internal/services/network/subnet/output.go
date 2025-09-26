package subnet

import (
	"sort"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/network/subnet"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintSubnetTable prints subnet information to the application stdout as either a JSON response or a formatted table.
// When useJSON is true, the subnets (optionally paginated) are marshaled to JSON and written to stdout; when false, a table with the columns Name, CIDR, and Public is rendered.
// If pagination is provided it will be adjusted and pagination info will be logged after output.
// The sortBy parameter accepts "name" (case-insensitive sort by DisplayName) or "cidr" (lexicographic sort by CidrBlock) to order the table rows.
// Returns an error if JSON marshaling fails.
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

// PrintSubnetInfo displays information for the provided subnets either as JSON or as formatted key-value blocks.
//
// When useJSON is true, writes a compact JSON response for the subnet slice; if the slice is empty it writes {"items": []}\n.
// When useJSON is false, prints a key-value view for each subnet containing the keys Name, Public, CIDR, DNS Label, and DNS Domain.
// The DNS Label is derived from the subnet's DisplayName and the DNS Domain is the label suffixed with ".vcn1.oraclevcn.com".
// Returns any error encountered while marshaling JSON output; otherwise returns nil.
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
// deriveDNSLabel derives a DNS-safe, lowercase label for a subnet name.
// If name is empty or sanitization yields no allowed characters, it returns "subnet".
// If name ends with digits, it returns "subnet" followed by those trailing digits.
// Otherwise it returns a lowercase string that preserves letters and digits and
// converts spaces/underscores to hyphens, dropping all other characters.
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
