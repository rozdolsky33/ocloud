package util

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
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
	coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
	coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
	coloredName := text.Colors{text.FgBlue}.Sprint(name)
	title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredName)

	return title
}

// SplitTextByMaxWidth splits a space-separated string into multiple lines
// SplitTextByMaxWidth splits a space-separated string into multiple lines where each line's length does not exceed 30 characters.
// If text is empty it returns a slice containing an empty string. Words are kept intact and joined into lines; a single-word input is returned unchanged.
// The returned slice contains the resulting lines.
func SplitTextByMaxWidth(text string) []string {

	if text == "" {
		return []string{""}
	}

	parts := strings.Fields(text)

	if len(parts) <= 1 {
		return []string{text}
	}

	result := make([]string, 0)
	currentLine := parts[0]

	for i := 1; i < len(parts); i++ {
		if len(currentLine)+len(parts[i])+1 > 30 {
			result = append(result, currentLine)
			currentLine = parts[i]
		} else {
			currentLine += " " + parts[i]
		}
	}

	result = append(result, currentLine)

	return result
}

// FormatBool returns a consistent string representation for booleans used in table outputs.
// FormatBool provides a consistent "Yes"/"No" string representation for boolean values used in table outputs.
// It returns "Yes" if b is true, "No" otherwise.
func FormatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
