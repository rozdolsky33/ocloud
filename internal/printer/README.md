# Printer Package

The `printer` package provides functionality to print data in tabular format using the [go-pretty](https://github.com/jedib0t/go-pretty) library. It offers a simple and flexible way to display structured data in the terminal with various formatting options.

## Overview

The `Printer` struct is the main component of this package, providing methods to print data in different formats:

- JSON output with proper indentation
- Key-value tables with ordered keys, titles, and colored values

## Installation

This package is part of the ocloud project and is already included in the codebase. To use it in your code, simply import it:

```
import "github.com/rozdolsky33/ocloud/internal/printer"
```

## Dependencies

The package relies on the following external dependencies:

- [github.com/jedib0t/go-pretty/v6](https://github.com/jedib0t/go-pretty) - For table rendering and styling

## Usage

### Creating a Printer

To create a new `Printer` that writes to a specific output, use the `New` function:

```
// For console output
p := printer.New(os.Stdout)

// For testing with a buffer
var buf bytes.Buffer
p := printer.New(&buf)
```

### Marshaling Data to JSON

To marshal data to JSON and write it to the printer's output:

```
// Create a data map
dataMap := map[string]string{
    "ID": "ocid1.instance.oc1.iad.123",
    "Name": "my-instance",
    "State": "RUNNING",
}

// Marshal to JSON and write to output
err := p.MarshalToJSON(dataMap)
if err != nil {
    // Handle error
}
```

This will produce output similar to:

```json
{
  "ID": "ocid1.instance.oc1.iad.123",
  "Name": "my-instance",
  "State": "RUNNING"
}
```

### Printing a Key-Value Table with Ordered Keys and Colored Values

To print a key-value table with ordered keys, a title, and colored values:

```
// Create a data map
instanceData := map[string]string{
    "ID": "ocid1.instance.oc1.iad.123",
    "Name": "my-instance",
    "State": "RUNNING",
    "Region": "us-ashburn-1",
    "Shape": "VM.Standard.E5.Flex",
}

// Define the order of keys
orderedKeys := []string{
    "ID",
    "Name",
    "State",
    "Region",
    "Shape",
}

// Print the table with a title and ordered keys
p.PrintKeyValues("Instance Details", instanceData, orderedKeys)
```

This will produce output similar to:

```
╭──────────────────────────────────────────────────────────────────╮
│                       Instance Details                           │
├─────┬──────────────────────────────────────────────────────────┤
│ KEY │ VALUE                                                    │
├─────┼──────────────────────────────────────────────────────────┤
│ ID  │ ocid1.instance.oc1.iad.123                               │
├─────┼──────────────────────────────────────────────────────────┤
│ Name│ my-instance                                              │
├─────┼──────────────────────────────────────────────────────────┤
│ State│ RUNNING                                                 │
├─────┼──────────────────────────────────────────────────────────┤
│ Region│ us-ashburn-1                                           │
├─────┼──────────────────────────────────────────────────────────┤
│ Shape│ VM.Standard.E5.Flex                                     │
╰─────┴──────────────────────────────────────────────────────────╯
```

In this example:
- The values are colored yellow for better visibility
- Each key-value pair is separated by a horizontal line
- The title is centered at the top of the table

## Real-World Example

Here's a real-world example from the ocloud project that shows how to use the `Printer` to display instance information:

```
// PrintInstancesTable displays instances in a formatted table or JSON format.
func PrintInstancesTable(instances []Instance, appCtx *app.ApplicationContext, pagination *PaginationInfo, useJSON bool) error {
    // Create a new printer that writes to the application's standard output.
    p := printer.New(appCtx.Stdout)

    // If JSON output is requested, use the printer to marshal the response.
    if useJSON {
        return marshalInstancesToJSON(p, instances, pagination)
    }

    // Handle the case where no instances are found.
    if len(instances) == 0 {
        fmt.Fprintln(appCtx.Stdout, "No instances found.")
        if pagination != nil && pagination.TotalCount > 0 {
            fmt.Fprintf(appCtx.Stdout, "Page %d is empty. Total records: %d\n", 
                pagination.CurrentPage, pagination.TotalCount)
            if pagination.CurrentPage > 1 {
                fmt.Fprintf(appCtx.Stdout, "Try a lower page number (e.g., --page %d)\n", 
                    pagination.CurrentPage-1)
            }
        }
        return nil
    }

    // Print each instance as a separate key-value table with a colored title.
    for _, instance := range instances {
        instanceData := map[string]string{
            "ID":         instance.ID,
            "AD":         instance.Placement.AvailabilityDomain,
            "FD":         instance.Placement.FaultDomain,
            "Region":     instance.Placement.Region,
            "Shape":      instance.Shape,
            "vCPUs":      fmt.Sprintf("%d", instance.Resources.VCPUs),
            "Created":    instance.CreatedAt.String(),
            "Subnet ID":  instance.SubnetID,
            "Name":       instance.Name,
            "Private IP": instance.IP,
            "Memory":     fmt.Sprintf("%d GB", int(instance.Resources.MemoryGB)),
            "State":      string(instance.State),
        }

        orderedKeys := []string{
            "ID", "AD", "FD", "Region", "Shape", "vCPUs",
            "Created", "Subnet ID", "Name", "Private IP", "Memory", "State",
        }

        // Create the colored title using components from the app context.
        coloredTenancy := text.Colors{text.FgMagenta}.Sprint(appCtx.TenancyName)
        coloredCompartment := text.Colors{text.FgCyan}.Sprint(appCtx.CompartmentName)
        coloredInstance := text.Colors{text.FgBlue}.Sprint(instance.Name)
        title := fmt.Sprintf("%s: %s: %s", coloredTenancy, coloredCompartment, coloredInstance)

        // Call the printer method to render the key-value table for this instance.
        p.PrintKeyValues(title, instanceData, orderedKeys)
    }

    return nil
}
```

## Styling Options

The `Printer` uses the `StyleRounded` style from the go-pretty library for key-value tables. This style provides a clean and readable output with rounded corners.

You can see all available styles in the [go-pretty documentation](https://github.com/jedib0t/go-pretty/blob/main/table/style.go).

## Color Options

The `PrintKeyValues` method applies yellow color to the values in the table:

- Values: Yellow (`text.FgYellow`)

This coloring helps to visually distinguish the values from the keys and makes the output more readable.

In the real-world example, additional colors are applied to the title components:

- Tenancy name: Magenta (`text.FgMagenta`)
- Compartment name: Cyan (`text.FgCyan`)
- Instance name: Blue (`text.FgBlue`)

## Customization

While the `Printer` provides sensible defaults, you can customize the output by modifying the table style, colors, and other properties directly in your code. For example:

```
// Create a custom table
t := table.NewWriter()
t.SetOutputMirror(myWriter)
t.SetStyle(table.StyleDouble) // Use double borders
t.Style().Title.Align = text.AlignLeft // Align title to the left
t.Style().Options.DrawBorder = false // Don't draw the outer border

// Add content to the table
t.AppendHeader(table.Row{"KEY", "VALUE"})
t.AppendRow(table.Row{"ID", "ocid1.instance.oc1.iad.123"})
t.AppendRow(table.Row{"Name", "my-instance"})

// Render the table
t.Render()
```

For more customization options, refer to the [go-pretty documentation](https://github.com/jedib0t/go-pretty).

## Conclusion

The `printer` package provides a simple and flexible way to display structured data in the terminal. It's particularly useful for displaying information about cloud resources, such as instances, in a readable and visually appealing format. The package supports both JSON output and tabular output, making it versatile for different use cases.
