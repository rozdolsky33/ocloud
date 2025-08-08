package bastion

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

func PrintBastionInfo(bastions []Bastion, appCtx *app.ApplicationContext, useJSON bool) error {

	// Create a new printer that writes to the application's standard output.
	p := printer.New(appCtx.Stdout)
	if useJSON {
		if len(bastions) == 0 {
			return p.MarshalToJSON(struct{}{})
		}
		return p.MarshalToJSON(bastions)
	}

	for _, b := range bastions {
		bastionInfo := map[string]string{
			"Name":           b.Name,
			"BastionType":    string(b.BastionType),
			"LifecycleState": string(b.LifecycleState),
			"TargetVcn":      b.TargetVcnName,
			"TargetSubnet":   b.TargetSubnetName,
		}
		// Define ordered Keys
		orderedKeys := []string{
			"Name", "BastionType", "LifecycleState", "TargetVcn", "TargetSubnet",
		}

		title := util.FormatColoredTitle(appCtx, b.Name)

		// Call the printer method to render the key-value table for this instance.
		p.PrintKeyValues(title, bastionInfo, orderedKeys)
	}

	return nil
}

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
