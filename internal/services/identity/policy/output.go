package policy

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintPolicyInfo prints the details of policies to the standard output or in JSON format.
// If pagination info is provided, it adjusts and logs it.
func PrintPolicyInfo(policies []Policy, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {

	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	// If JSON output is requested, use the printer to marshal the response.
	if useJSON {
		return util.MarshalDataToJSONResponse[Policy](p, policies, pagination)
	}

	if util.ValidateAndReportEmpty(policies, pagination, appCtx.Stdout) {
		return nil
	}

	// Print each policy as a separate key-value table with a colored title,
	for _, policy := range policies {
		policyData := map[string]string{
			"Name":        policy.Name,
			"ID":          policy.ID,
			"Description": policy.Description,
		}

		// Define ordered keys
		orderedKeys := []string{
			"Name", "ID", "Description",
		}

		// Create the colored title using components from the app context
		title := util.FormatColoredTitle(appCtx, policy.Name)

		// Call the printer method to render the key-value from the app context.
		p.PrintKeyValues(title, policyData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
