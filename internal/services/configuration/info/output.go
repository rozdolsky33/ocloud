package info

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintMappingsFile displays tenancy mapping information in a formatted table or JSON format.
// It takes a slice of MappingsFile, the application context, and a boolean indicating whether to use JSON format.
// Returns an error if the display operation fails.
func PrintMappingsFile(mappings []appConfig.MappingsFile, useJSON bool) error {

	p := printer.New(os.Stdout)

	if useJSON {
		if len(mappings) == 0 {
			return p.MarshalToJSON(struct{}{})
		}
		return p.MarshalToJSON(mappings)
	}

	if util.ValidateAndReportEmpty(mappings, nil, os.Stdout) {
		return nil
	}

	// Group mappings by realm
	realmGroups := groupMappingsByRealm(mappings)

	// headers for the table
	headers := []string{"ENVIRONMENT", "TENANCY", "COMPARTMENTS", "REGIONS"}

	// For each realm, create and display a separate table
	for realm, realmMappings := range realmGroups {
		// Convert mappings to rows for the table, handling long compartment names and regions
		rows := make([][]string, 0, len(realmMappings))
		for _, mapping := range realmMappings {
			compart := strings.Join(mapping.Compartments, " ")
			reg := strings.Join(mapping.Regions, " ")
			compartments := util.SplitTextByMaxWidth(compart)
			regions := util.SplitTextByMaxWidth(reg)

			// Create the first row with all columns
			firstRow := []string{
				mapping.Environment,
				mapping.Tenancy,
				compartments[0],
				regions[0],
			}
			rows = append(rows, firstRow)

			maxAdditionalRows := len(compartments) - 1
			if len(regions)-1 > maxAdditionalRows {
				maxAdditionalRows = len(regions) - 1
			}

			// Create additional rows for compartments and regions if needed
			for i := 0; i < maxAdditionalRows; i++ {
				compartment := ""
				if i+1 < len(compartments) {
					compartment = compartments[i+1]
				}

				region := ""
				if i+1 < len(regions) {
					region = regions[i+1]
				}

				additionalRow := []string{
					"",
					"",
					compartment,
					region,
				}
				rows = append(rows, additionalRow)
			}
		}

		coloredTitle := text.Colors{text.FgMagenta}.Sprint(fmt.Sprintf("Tenancy Mapping Information - Realm: %s", realm))
		p.PrintTable(coloredTitle, headers, rows)
	}

	return nil
}

// groupMappingsByRealm groups mappings by their realm.
// It returns a map where the key is the realm and the value is a slice of mappings for that realm.
func groupMappingsByRealm(mappings []appConfig.MappingsFile) map[string][]appConfig.MappingsFile {
	realmGroups := make(map[string][]appConfig.MappingsFile)

	for _, mapping := range mappings {
		realmGroups[mapping.Realm] = append(realmGroups[mapping.Realm], mapping)
	}

	return realmGroups
}
