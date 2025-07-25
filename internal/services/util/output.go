package util

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"strings"
)

// MarshalDataToJSONResponse now accepts a printer and returns an error.
func MarshalDataToJSONResponse[T any](p *printer.Printer, items []T, pagination *PaginationInfo) error {
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

// SplitTextByMaxWidth splits a space-separated string into multiple lines
// to ensure they are all visible in the table output with a maximum width per line.
func SplitTextByMaxWidth(text string) []string {

	if text == "" {
		return []string{""}
	}

	parts := strings.Fields(text)

	// If there's only one part, or it's short enough, return it as is
	if len(parts) <= 1 {
		return []string{text}
	}

	result := make([]string, 0)
	currentLine := parts[0]

	for i := 1; i < len(parts); i++ {
		// If adding the next part makes the line too long, start a new line
		if len(currentLine)+len(parts[i])+1 > 30 { // 30 is a reasonable width for the column
			result = append(result, currentLine)
			currentLine = parts[i]
		} else {
			currentLine += " " + parts[i]
		}
	}

	result = append(result, currentLine)

	return result
}
