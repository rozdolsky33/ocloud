package util

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
)

// MarshalDataToJSON now accepts a printer and returns an error.
func MarshalDataToJSON[T any](p *printer.Printer, items []T, pagination *PaginationInfo) error {
	response := JSONResponse[T]{
		Items:      items,
		Pagination: pagination,
	}
	return p.MarshalToJSON(response)
}

// FormatColoredTitle builds a colorized title string with tenancy, compartment, and cluster.
func FormatColoredTitle(appCtx *app.ApplicationContext, name string) string {
	// Create the colored title using components from the app context.
	coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
	coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
	coloredName := text.Colors{text.FgBlue}.Sprint(name)
	title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredName)

	return title
}
