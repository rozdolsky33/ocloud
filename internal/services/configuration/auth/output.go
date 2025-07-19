package auth

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"strings"
)

// AuthenticateWithOCI handles the authentication process with OCI.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
// If envOnly is true, it skips the interactive prompts and only outputs the environment variables.
// If filter is not empty, it filters the regions by prefix.
func AuthenticateWithOCI(appCtx *app.ApplicationContext, envOnly bool, filter string) error {
	logger.LogWithLevel(appCtx.Logger, 1, "Authenticating with OCI", "envOnly", envOnly)

	// Create a new service
	service := NewService(appCtx)

	var result *AuthenticationResult
	var err error

	if envOnly {
		// In env-only mode, skip interactive prompts and authentication
		// Just get the current environment variables
		result, err = service.GetCurrentEnvironment()
		if err != nil {
			return fmt.Errorf("getting current environment: %w", err)
		}
	} else {
		// Prompt for profile selection
		profile, err := service.PromptForProfile()
		if err != nil {
			return fmt.Errorf("selecting profile: %w", err)
		}

		// Display regions in a table
		regions := service.GetOCIRegions()
		if err := DisplayRegionsTable(regions, appCtx, filter); err != nil {
			return fmt.Errorf("displaying regions: %w", err)
		}

		// Prompt for region selection
		region, err := service.PromptForRegion()
		if err != nil {
			return fmt.Errorf("selecting region: %w", err)
		}
		fmt.Printf("Using region: %s\n", region)

		// Authenticate with OCI
		result, err = service.Authenticate(profile, region)
		if err != nil {
			return fmt.Errorf("authenticating with OCI: %w", err)
		}
	}

	// Print export for compartment and tenancy
	fmt.Printf("export OCI_COMPARTMENT=%s\n", result.TenancyID)
	fmt.Printf("export OCI_CLI_TENANCY=%s\n", result.TenancyID)

	// Print export for tenancy name if available
	if result.TenancyName != "" {
		fmt.Printf("export OCI_TENANCY_NAME=%s\n", result.TenancyName)
	}

	if !envOnly {
		fmt.Println("✅ Authentication complete. Run `eval $(ocloud config auth -e)` to set your env.")
	}

	return nil
}

// DisplayRegionsTable displays the available OCI regions in a table format.
// If filter is not empty, it filters the regions by prefix.
func DisplayRegionsTable(regions []RegionInfo, appCtx *app.ApplicationContext, filter string) error {
	// Create a new printer
	p := printer.New(appCtx.Stdout)

	// Group regions by their prefix (e.g., us, eu, ap)
	regionGroups := groupRegionsByPrefix(regions)

	// Filter regions by prefix if filter is provided
	if filter != "" {
		// Convert filter to lowercase for case-insensitive comparison
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
		// Create a title for the group with magenta color
		regionTitle := getRegionGroupTitle(prefix)
		groupTitle := text.Colors{text.FgMagenta}.Sprint(fmt.Sprintf("%s", regionTitle))

		// Create rows for the regions table
		var rows [][]string

		// Add an empty row for spacing
		rows = append(rows, []string{""})

		// Add the regions in rows with multiple regions per row
		var currentRegions []string

		for i, region := range prefixRegions {
			// Format each region as "ID: region-name" with color
			regionName := text.Colors{text.FgGreen}.Sprint(region.Name)
			regionID := text.Colors{text.FgRed}.Sprint(region.ID)
			formattedRegion := fmt.Sprintf("%s: %s", regionID, regionName)
			currentRegions = append(currentRegions, formattedRegion)

			// Start a new row after every 3 regions or at the end
			if (i+1)%5 == 0 || i == len(prefixRegions)-1 {
				rows = append(rows, []string{strings.Join(currentRegions, "  ")})
				currentRegions = nil
			}
		}

		// Call the printer method to render the table with the region title and regions
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
