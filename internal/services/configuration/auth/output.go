package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

// DisplayRegionsTable displays the available OCI regions in a table format.
func DisplayRegionsTable(regions []RegionInfo, filter string) error {

	p := printer.New(os.Stdout)

	// Group regions by their prefix (e.g., us, eu, ap)
	regionGroups := groupRegionsByPrefix(regions)

	// Filter regions by prefix if filter is provided
	if filter != "" {
		filter = strings.ToLower(filter)
		filteredGroups := make(map[string][]RegionInfo)
		for prefix, prefixRegions := range regionGroups {
			if strings.HasPrefix(strings.ToLower(prefix), filter) {
				filteredGroups[prefix] = prefixRegions
			}
		}

		regionGroups = filteredGroups
	}

	// Process each region group
	for prefix, prefixRegions := range regionGroups {
		regionTitle := getRegionGroupTitle(prefix)
		groupTitle := text.Colors{text.FgMagenta}.Sprint(fmt.Sprintf("%s", regionTitle))

		var rows [][]string
		rows = append(rows, []string{""})
		var currentRegions []string

		for i, region := range prefixRegions {
			regionName := text.Colors{text.FgGreen}.Sprint(region.Name)
			regionID := text.Colors{text.FgRed}.Sprint(region.ID)
			formattedRegion := fmt.Sprintf("%s: %s", regionID, regionName)
			currentRegions = append(currentRegions, formattedRegion)
			if (i+1)%4 == 0 || i == len(prefixRegions)-1 {
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

// PrintExportVariable prints the environment variables in a centered table with color.
func PrintExportVariable(profile, tenancyName, compartment string) error {
	logger.LogWithLevel(logger.Logger, 3, "Printing export variables", "profile", profile, "tenancyName", tenancyName, "compartment", compartment)

	// Create a map of environment variables to export
	exportVars := make(map[string]string)

	if profile != "" {
		exportVars[flags.EnvKeyProfile] = profile
		logger.LogWithLevel(logger.Logger, 3, "Added profile to export variables", "profile", profile)
	}

	if tenancyName != "" {
		exportVars[flags.EnvKeyTenancyName] = tenancyName
		logger.LogWithLevel(logger.Logger, 3, "Added tenancy name to export variables", "tenancyName", tenancyName)
	}

	if compartment != "" {
		exportVars[flags.EnvKeyCompartment] = compartment
		logger.LogWithLevel(logger.Logger, 3, "Added compartment to export variables", "compartment", compartment)
	}

	// Create a printer and print the export variables in a table
	p := printer.New(os.Stdout)
	title := "Export Variable"
	message := "ENVIRONMENT VARIABLES"
	p.ResultTable(title, message, exportVars)

	logger.LogWithLevel(logger.Logger, 3, "Printed export variables in table")

	fmt.Println("\nTo persist your selection, export the following environment variables in your shell")

	return nil
}
