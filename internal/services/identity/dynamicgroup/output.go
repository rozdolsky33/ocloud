package dynamicgroup

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintDynamicGroupsTable displays a table of dynamic groups.
func PrintDynamicGroupsTable(dynamicGroups []DynamicGroup, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}

	if useJSON {
		if len(dynamicGroups) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
		return util.MarshalDataToJSONResponse[DynamicGroup](p, dynamicGroups, pagination)
	}

	if util.ValidateAndReportEmpty(dynamicGroups, pagination, appCtx.Stdout) {
		return nil
	}

	headers := []string{"Name", "ID", "State"}

	rows := make([][]string, len(dynamicGroups))
	for i, dg := range dynamicGroups {
		rows[i] = []string{
			dg.Name,
			dg.OCID,
			dg.LifecycleState,
		}
	}

	title := util.FormatColoredTitle(appCtx, "Dynamic Groups")
	p.PrintTable(title, headers, rows)

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}

// PrintDynamicGroupInfo displays a detailed view of a dynamic group.
func PrintDynamicGroupInfo(dg *DynamicGroup, appCtx *app.ApplicationContext, useJSON bool) error {
	p := printer.New(appCtx.Stdout)

	if useJSON {
		return p.MarshalToJSON(dg)
	}
	dgData := map[string]string{
		"Name":          dg.Name,
		"ID":            dg.OCID,
		"Description":   dg.Description,
		"Matching Rule": dg.MatchingRule,
		"State":         dg.LifecycleState,
		"Created":       dg.TimeCreated.String(),
	}
	orderedKeys := []string{
		"Name", "ID", "Description", "Matching Rule", "State", "Created",
	}

	title := util.FormatColoredTitle(appCtx, dg.Name)

	p.PrintKeyValues(title, dgData, orderedKeys)

	return nil
}

// PrintDynamicGroupsInfo displays information about a list of dynamic groups.
func PrintDynamicGroupsInfo(dgs []DynamicGroup, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool) error {
	p := printer.New(appCtx.Stdout)
	if pagination != nil {
		util.AdjustPaginationInfo(pagination)
	}
	if useJSON {
		if len(dgs) == 0 && pagination == nil {
			return p.MarshalToJSON(struct{}{})
		}
		return util.MarshalDataToJSONResponse[DynamicGroup](p, dgs, pagination)
	}

	if util.ValidateAndReportEmpty(dgs, pagination, appCtx.Stdout) {
		return nil
	}

	for _, dg := range dgs {
		dgData := map[string]string{
			"Name":          dg.Name,
			"ID":            dg.OCID,
			"Description":   dg.Description,
			"Matching Rule": dg.MatchingRule,
			"State":         dg.LifecycleState,
		}
		orderedKeys := []string{
			"Name", "ID", "Description", "Matching Rule", "State",
		}

		title := util.FormatColoredTitle(appCtx, dg.Name)

		p.PrintKeyValues(title, dgData, orderedKeys)
	}

	util.LogPaginationInfo(pagination, appCtx)
	return nil
}
