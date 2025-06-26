package compartment

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintCompartmentsInfo(compartments []Compartment, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {

	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return util.MarshalDataToJSONResponse[Compartment](p, compartments, pagination)
	}

	if util.ValidateAndReportEmpty(compartments, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each Compartment as a separate key-value table with a colored title.
	for _, compartment := range compartments {
		compartmentData := map[string]string{
			"Name":        compartment.Name,
			"ID":          compartment.ID,
			"Description": compartment.Description,
		}
		// Define ordered keys
		orderedKeys := []string{
			"Name", "ID", "Description",
		}

		title := util.FormatColoredTitle(appCtx, compartment.Name)

		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, compartmentData, orderedKeys)
	}
	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
