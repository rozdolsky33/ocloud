package auth

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"strings"
)

// DisplayRegionsTable displays the available OCI regions in a table format.
// If the filter is not empty, it filters the regions by prefix.
func DisplayRegionsTable(regions []RegionInfo, appCtx *app.ApplicationContext, filter string) error {
	logger := appCtx.Logger
	logger.V(3).Info("Displaying regions table", "totalRegions", len(regions), "filter", filter)

	p := printer.New(appCtx.Stdout)

	// Group regions by their prefix (e.g., us, eu, ap)
	regionGroups := groupRegionsByPrefix(regions)
	logger.V(3).Info("Grouped regions by prefix", "groupCount", len(regionGroups))

	// Filter regions by prefix if filter is provided
	if filter != "" {
		logger.V(3).Info("Filtering regions by prefix", "filter", filter)
		filter = strings.ToLower(filter)
		// Create a new map with only the filtered regions
		filteredGroups := make(map[string][]RegionInfo)
		for prefix, prefixRegions := range regionGroups {
			if strings.HasPrefix(strings.ToLower(prefix), filter) {
				filteredGroups[prefix] = prefixRegions
				logger.V(3).Info("Including region group in filter", "prefix", prefix, "regionCount", len(prefixRegions))
			}
		}

		// Replace the original map with the filtered one
		regionGroups = filteredGroups
		logger.V(3).Info("After filtering", "groupCount", len(regionGroups))
	}

	// Process each region group
	for prefix, prefixRegions := range regionGroups {
		regionTitle := getRegionGroupTitle(prefix)
		groupTitle := text.Colors{text.FgMagenta}.Sprint(fmt.Sprintf("%s", regionTitle))
		logger.V(3).Info("Processing region group", "prefix", prefix, "title", regionTitle, "regionCount", len(prefixRegions))

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

	logger.V(3).Info("Finished displaying regions table")
	return nil
}

// groupRegionsByPrefix groups regions by their prefix (e.g., us, eu, ap).
func groupRegionsByPrefix(regions []RegionInfo) map[string][]RegionInfo {
	// Use the package-level logger since this is not a method
	logger.LogWithLevel(logger.Logger, 3, "Grouping regions by prefix", "regionCount", len(regions))

	regionGroups := make(map[string][]RegionInfo)

	for _, region := range regions {
		// Extract the prefix (e.g., "us" from "us-ashburn-1")
		parts := strings.Split(region.Name, "-")
		if len(parts) > 0 {
			prefix := parts[0]
			regionGroups[prefix] = append(regionGroups[prefix], region)
		}
	}

	// Log the results
	for prefix, regions := range regionGroups {
		logger.LogWithLevel(logger.Logger, 3, "Region group", "prefix", prefix, "count", len(regions))
	}

	return regionGroups
}

// getRegionGroupTitle returns a human-readable title for a region group.
func getRegionGroupTitle(prefix string) string {
	// Use the package-level logger since this is not a method
	logger.LogWithLevel(logger.Logger, 3, "Getting region group title", "prefix", prefix)

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
		logger.LogWithLevel(logger.Logger, 3, "Found title for prefix", "prefix", prefix, "title", title)
		return title
	}

	logger.LogWithLevel(logger.Logger, 3, "No title found for prefix, using prefix as title", "prefix", prefix)
	return prefix
}

// PrintExportVariable prints the environment variables with color and proper spacing.
// If tenancyName and compartment are provided, they are included in the output.
func PrintExportVariable(tenancyName, compartment string) error {
	logger.LogWithLevel(logger.Logger, 3, "Printing export variables", "tenancyName", tenancyName, "compartment", compartment)

	// Format the export statements with color
	tenancyNameVar := text.Colors{text.FgYellow}.Sprint("export OCI_TENANCY_NAME=")
	compartmentVar := text.Colors{text.FgYellow}.Sprint("export OCI_COMPARTMENT=")

	// Add the values if provided
	if tenancyName != "" {
		tenancyNameVar += fmt.Sprintf("%q", tenancyName)
		logger.LogWithLevel(logger.Logger, 3, "Added tenancy name to export variable", "tenancyName", tenancyName)
	}

	if compartment != "" {
		compartmentVar += fmt.Sprintf("%q", compartment)
		logger.LogWithLevel(logger.Logger, 3, "Added compartment to export variable", "compartment", compartment)
	}

	// Print with proper spacing
	fmt.Println(tenancyNameVar)
	fmt.Println(compartmentVar)
	logger.LogWithLevel(logger.Logger, 3, "Printed export variables")

	return nil
}
