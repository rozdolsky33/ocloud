package auth

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"strings"
)

// DisplayRegionsTable displays the available OCI regions in a table format.
// If the filter is not empty, it filters the regions by prefix.
func DisplayRegionsTable(regions []RegionInfo, appCtx *app.ApplicationContext, filter string) error {
	p := printer.New(appCtx.Stdout)

	// Group regions by their prefix (e.g., us, eu, ap)
	regionGroups := groupRegionsByPrefix(regions)

	// Filter regions by prefix if filter is provided
	if filter != "" {
		filter = strings.ToLower(filter)
		// Create a new map with only the filtered regions
		filteredGroups := make(map[string][]RegionInfo)
		for prefix, prefixRegions := range regionGroups {
			if strings.HasPrefix(strings.ToLower(prefix), filter) {
				filteredGroups[prefix] = prefixRegions
			}
		}

		// Replace the original map with the filtered one
		regionGroups = filteredGroups
	}

	// Process each region group
	for prefix, prefixRegions := range regionGroups {
		regionTitle := getRegionGroupTitle(prefix)
		groupTitle := text.Colors{text.FgMagenta}.Sprint(fmt.Sprintf("%s", regionTitle))

		// Create rows for the region table
		var rows [][]string

		rows = append(rows, []string{""})

		var currentRegions []string

		for i, region := range prefixRegions {
			regionName := text.Colors{text.FgGreen}.Sprint(region.Name)
			regionID := text.Colors{text.FgRed}.Sprint(region.ID)
			formattedRegion := fmt.Sprintf("%s: %s", regionID, regionName)
			currentRegions = append(currentRegions, formattedRegion)
			if (i+1)%5 == 0 || i == len(prefixRegions)-1 {
				rows = append(rows, []string{strings.Join(currentRegions, "  ")})
				currentRegions = nil
			}
		}
		p.PrintTable(groupTitle, []string{"Available OCI Regions"}, rows)
	}

	return nil
}

// groupRegionsByPrefix groups regions by their prefix (e.g., us, eu, ap).
func groupRegionsByPrefix(regions []RegionInfo) map[string][]RegionInfo {
	regionGroups := make(map[string][]RegionInfo)

	for _, region := range regions {
		// Extract the prefix (e.g., "us" from "us-ashburn-1")
		parts := strings.Split(region.Name, "-")
		if len(parts) > 0 {
			prefix := parts[0]
			regionGroups[prefix] = append(regionGroups[prefix], region)
		}
	}

	return regionGroups
}

// getRegionGroupTitle returns a human-readable title for a region group.
func getRegionGroupTitle(prefix string) string {
	titles := map[string]string{
		"af": "Africa",
		"ap": "Asia Pacific",
		"ca": "Canada",
		"eu": "Europe",
		"il": "Israel",
		"me": "Middle East",
		"mx": "Mexico",
		"sa": "South America",
		"uk": "United Kingdom",
		"us": "United States",
	}

	if title, ok := titles[prefix]; ok {
		return title
	}
	return prefix
}

func PrintExportVariable() error {
	tenancyName := text.Colors{text.FgYellow}.Sprint("export OCI_TENANCY_NAME=")
	compartment := text.Colors{text.FgYellow}.Sprint("export OCI_COMPARTMENT=", "\n")
	fmt.Printf(tenancyName + "\n")
	fmt.Printf(compartment)
	return nil
}
