package bastion

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintBastions prints a list of bastions in a tabular format
func PrintBastions(bastions []Bastion) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\t")
	fmt.Fprintln(w, "----\t----\t")

	for _, b := range bastions {
		fmt.Fprintf(w, "%s\t%s\t\n", b.ID, b.Name)
	}

	w.Flush()
}
