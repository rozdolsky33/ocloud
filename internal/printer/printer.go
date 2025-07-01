package printer

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"io"
)

// Printer handles formatting and writing output to a designated writer.
type Printer struct {
	out io.Writer
}

// New creates a new Printer that writes to the provided io.Writer.
// For console output, use os.Stdout. For testing, use bytes.Buffer.
func New(out io.Writer) *Printer {
	return &Printer{out: out}
}

// MarshalToJSON marshals data to JSON and writes it to the printer's output.
func (p *Printer) MarshalToJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	// Use fmt.Fprintln to write to the specific io.Writer.
	_, err = fmt.Fprintln(p.out, string(jsonData))
	return err
}

// PrintKeyValues renders a table from a map, with ordered keys, a title, and colored values.
// This new method encapsulates all the logic for your detailed instance view.
func (p *Printer) PrintKeyValues(title string, data map[string]string, keys []string) {
	t := table.NewWriter()
	t.SetOutputMirror(p.out)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter
	t.SetTitle(title)

	header := table.Row{"KEY", "VALUE"}
	t.AppendHeader(header)

	for i, key := range keys {
		if value, ok := data[key]; ok {
			if i > 0 {
				t.AppendSeparator()
			}
			// Color the value text yellow, just like in the original.
			coloredValue := text.Colors{text.FgYellow}.Sprint(value)
			t.AppendRow(table.Row{key, coloredValue})
		}
	}

	t.Render()
}

// PrintTable renders a table with the given headers and rows.
// This method is used for displaying data in a tabular format.
func (p *Printer) PrintTable(title string, headers []string, rows [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(p.out)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter
	t.SetTitle(title)

	// Convert headers to table.Row

	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = text.Colors{text.FgHiYellow}.Sprint(h)
	}

	t.AppendHeader(headerRow)

	// Add rows to the table
	for _, row := range rows {
		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = cell
		}
		t.AppendRow(tableRow)
	}

	t.Render()
}
