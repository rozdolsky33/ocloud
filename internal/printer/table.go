package printer

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// TablePrinter provides functionality to print data in tabular format
type TablePrinter struct {
	Title string
}

// NewTablePrinter creates a new table printer
func NewTablePrinter(title string) *TablePrinter {
	return &TablePrinter{
		Title: title,
	}
}

// PrintKeyValueTable prints data in a key-value table format
func (p *TablePrinter) PrintKeyValueTable(data []map[string]string) {
	if len(data) == 0 {
		fmt.Println("No data to display.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	// Set header
	header := table.Row{"KEY", "VALUE"}
	t.AppendHeader(header)

	// For each data item, create a section in the table
	for _, item := range data {
		// Add all key-value pairs
		for key, value := range item {
			t.AppendRow(table.Row{key, value})
		}
	}

	// Render the table
	t.Render()
}

// PrintKeyValueTableWithTitle prints data in a key-value table format with a custom title
func (p *TablePrinter) PrintKeyValueTableWithTitle(title string, data map[string]string) {
	if len(data) == 0 {
		fmt.Println("No data to display.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter
	t.SetTitle(title)

	// Set header
	header := table.Row{"KEY", "VALUE"}
	t.AppendHeader(header)

	// Add all key-value pairs
	for key, value := range data {
		t.AppendRow(table.Row{key, value})
	}

	// Render the table
	t.Render()
}

// PrintKeyValueTableWithTitleOrdered prints data in a key-value table format with a custom title
// and in the order specified by the key slice
func (p *TablePrinter) PrintKeyValueTableWithTitleOrdered(tenancyName, compartmentName, instanceName string, data map[string]string, keys []string) {
	if len(data) == 0 {
		fmt.Println("No data to display.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter

	// Color the different parts of the title
	coloredTenancy := text.Colors{text.FgMagenta}.Sprint(tenancyName)
	coloredCompartment := text.Colors{text.FgCyan}.Sprint(compartmentName)
	coloredInstance := text.Colors{text.FgBlue}.Sprint(instanceName)

	// Combine the colored parts into a single title
	title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredInstance)
	t.SetTitle(title)

	// Set header
	header := table.Row{"KEY", "VALUE"}
	t.AppendHeader(header)

	// Add key-value pairs in the specified order with separators and yellow values
	for i, key := range keys {
		if value, ok := data[key]; ok {
			// Add a separator before each row except the first one
			if i > 0 {
				t.AppendSeparator()
			}

			// Color the value text yellow but keep the key as is
			coloredValue := text.Colors{text.FgYellow}.Sprint(value)
			t.AppendRow(table.Row{key, coloredValue})
		}
	}

	// Render the table
	t.Render()
}

// PrintRowsTable prints data in a row table format
func (p *TablePrinter) PrintRowsTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Println("No data to display.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Title.Align = text.AlignCenter

	if p.Title != "" {
		t.SetTitle(p.Title)
	}

	// Set header
	headerRow := table.Row{}
	for _, header := range headers {
		headerRow = append(headerRow, header)
	}
	t.AppendHeader(headerRow)

	// Add each row
	for _, row := range rows {
		tableRow := table.Row{}
		for _, cell := range row {
			tableRow = append(tableRow, cell)
		}
		t.AppendRow(tableRow)
	}

	// Render the table
	t.Render()
}
