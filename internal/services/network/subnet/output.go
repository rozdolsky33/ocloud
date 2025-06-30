package subnet

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintSubnetInfo(subnets []Subnet, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {

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

	return nil
}
