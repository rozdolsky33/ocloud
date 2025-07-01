package printer

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
	"unicode/utf8"
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

// -----------------------------------------------------------------------------
// Utility helpers
// -----------------------------------------------------------------------------

// getTerminalWidth returns the current terminal width. If the writer is not a
// file descriptor (e.g., in tests) or the call fails, it falls back to 80 cols.
func (p *Printer) getTerminalWidth() int {
	if f, ok := p.out.(*os.File); ok {
		if w, _, err := term.GetSize(int(f.Fd())); err == nil {
			return w
		}
	}
	return 80 // sensible default
}

// truncate shortens a string to max runes, appending an ellipsis when needed.
func truncate(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	if max <= 3 {
		return string(r[:max])
	}
	return string(r[:max-3]) + "..."
}

// -----------------------------------------------------------------------------
// JSON output helpers
// -----------------------------------------------------------------------------

// MarshalToJSON marshals data to JSON and writes it to the printer's output.
func (p *Printer) MarshalToJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	_, err = fmt.Fprintln(p.out, string(jsonData))
	return err
}

// -----------------------------------------------------------------------------
// Key/Value table
// -----------------------------------------------------------------------------

// PrintKeyValues renders a table from a map, with ordered keys, a title, and
// colored values.  Useful for detailed instance/subnet views.
func (p *Printer) PrintKeyValues(title string, data map[string]string, keys []string) {
	t := table.NewWriter()
	t.SetOutputMirror(p.out)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter
	t.SetTitle(title)

	t.AppendHeader(table.Row{"KEY", "VALUE"})

	for i, key := range keys {
		if value, ok := data[key]; ok {
			if i > 0 {
				t.AppendSeparator()
			}
			coloredValue := text.Colors{text.FgYellow}.Sprint(value)
			t.AppendRow(table.Row{key, coloredValue})
		}
	}

	t.Render()
}

// -----------------------------------------------------------------------------
// Responsive multi‑column table
// -----------------------------------------------------------------------------

// PrintTable renders a table with the given headers and rows, automatically
// adapting column widths based on the current terminal size.
func (p *Printer) PrintTable(title string, headers []string, rows [][]string) {
	termWidth := p.getTerminalWidth()

	// Calculate a reasonable max width per column.
	// Rough formula: subtract borders/padding (≈3 chars per col), then divide.
	pad := (len(headers) + 1) * 3
	maxPerCol := (termWidth - pad) / len(headers)
	if maxPerCol < 10 {
		maxPerCol = 10 // never let columns get absurdly narrow
	}

	// Set up the table writer
	t := table.NewWriter()
	t.SetOutputMirror(p.out)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter
	t.SetTitle(title)

	// Build header row and column configs simultaneously
	headerRow := make(table.Row, len(headers))
	colConfigs := make([]table.ColumnConfig, len(headers))

	for i, h := range headers {
		headerRow[i] = text.Colors{text.FgHiYellow}.Sprint(h)

		// Capture i and h for the transformer closure
		idx := i
		colConfigs[i] = table.ColumnConfig{
			Number:   i + 1, // 1‑based index
			WidthMax: maxPerCol,
			Transformer: func(val interface{}) string {
				return truncate(fmt.Sprint(val), maxPerCol)
			},
			// Optional: align numeric/text differently if desired
			// Align: table.AlignLeft,
		}

		// Special case: align CIDR and IP columns to the center for readability
		if h == "CIDR" || strings.Contains(strings.ToLower(h), "ip") {
			colConfigs[idx].Align = text.AlignCenter
		}
	}

	t.AppendHeader(headerRow)
	t.SetColumnConfigs(colConfigs)

	// Add rows
	for _, row := range rows {
		tblRow := make(table.Row, len(row))
		for i, cell := range row {
			tblRow[i] = cell
		}
		t.AppendRow(tblRow)
	}

	t.Render()
}
