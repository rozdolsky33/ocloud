package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintAutonomousDbInfo displays instances in a formatted table or JSON format.
// It now returns an error to allow for proper error handling by the caller.
func PrintAutonomousDbInfo(databases []AutonomousDatabase, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)

	// Adjust the pagination information if available
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return util.MarshalDataToJSON[AutonomousDatabase](p, databases, pagination)
	}

	if util.ValidateAndReportEmpty(databases, pagination, appCtx.Stdout) {
		return nil
	}
	// Print each Compartment as a separate key-value table with a colored title.
	for _, database := range databases {
		databaseData := map[string]string{
			"Private IP":       database.PrivateEndpointIp,
			"ID":               database.ID,
			"Private Endpoint": database.PrivateEndpoint,
			"High":             database.ConnectionStrings["HIGH"],
			"Medium":           database.ConnectionStrings["MEDIUM"],
			"Low":              database.ConnectionStrings["LOW"],
		}
		// Define ordered Keys
		orderedKeys := []string{
			"Private IP", "ID", "Private Endpoint", "High", "Medium", "Low",
		}

		title := util.FormatColoredTitle(appCtx, database.Name)

		// Call the printer method to tender the key-value table for this instance
		p.PrintKeyValues(title, databaseData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
