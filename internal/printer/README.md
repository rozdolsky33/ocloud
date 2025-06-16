# TablePrinter Package

The `printer` package provides functionality to print data in tabular format using the [go-pretty](https://github.com/jedib0t/go-pretty) library. It offers a simple and flexible way to display structured data in the terminal with various formatting options.

## Overview

The `TablePrinter` struct is the main component of this package, providing methods to print data in different tabular formats:

- Key-value tables
- Key-value tables with custom titles
- Key-value tables with ordered keys and colored components
- Row-based tables with headers

## Installation

This package is part of the ocloud project and is already included in the codebase. To use it in your code, simply import it:

```
import "github.com/rozdolsky33/ocloud/internal/printer"
```

## Dependencies

The package relies on the following external dependencies:

- [github.com/jedib0t/go-pretty/v6](https://github.com/jedib0t/go-pretty) - For table rendering and styling
- [github.com/rozdolsky33/ocloud/internal/app](https://github.com/rozdolsky33/ocloud) - For application context

## Usage

### Creating a TablePrinter

To create a new `TablePrinter`, use the `NewTablePrinter` function:

```
tablePrinter := printer.NewTablePrinter("My Title")
```

### Printing a Simple Key-Value Table

To print a simple key-value table:

```
data := []map[string]string{
    {
        "ID": "ocid1.instance.oc1.iad.123",
        "Name": "my-instance",
        "State": "RUNNING",
    },
}
tablePrinter.PrintKeyValueTable(data)
```

This will produce output similar to:

```
╭─────┬──────────────────────────────────────╮
│ KEY │ VALUE                                │
├─────┼──────────────────────────────────────┤
│ ID  │ ocid1.instance.oc1.iad.123           │
│ Name│ my-instance                          │
│ State│ RUNNING                             │
╰─────┴──────────────────────────────────────╯
```

### Printing a Key-Value Table with a Custom Title

To print a key-value table with a custom title:

```
data := map[string]string{
    "ID": "ocid1.instance.oc1.iad.123",
    "Name": "my-instance",
    "State": "RUNNING",
}
tablePrinter.PrintKeyValueTableWithTitle("Instance Details", data)
```

This will produce output similar to:

```
╭──────────────────────────────────────────────────────────────────╮
│                       Instance Details                           │
├─────┬──────────────────────────────────────────────────────────┤
│ KEY │ VALUE                                                    │
├─────┼──────────────────────────────────────────────────────────┤
│ ID  │ ocid1.instance.oc1.iad.123                               │
│ Name│ my-instance                                              │
│ State│ RUNNING                                                 │
╰─────┴──────────────────────────────────────────────────────────╯
```

### Printing a Key-Value Table with Ordered Keys and Colored Components

To print a key-value table with ordered keys and colored components:

```
// Create application context
appCtx := &app.AppContext{
    TenancyName: "MyTenancy",
    CompartmentName: "MyCompartment",
}

// Create data map
data := map[string]string{
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

// Print the table with ordered keys and colored title components
tablePrinter.PrintKeyValueTableWithTitleOrdered(appCtx, "my-instance", data, orderedKeys)
```

This will produce output similar to:

```
╭──────────────────────────────────────────────────────────────────────────────────────────────────╮
│ MyTenancy: MyCompartment: my-instance                                                            │
├────────────┬─────────────────────────────────────────────────────────────────────────────────────┤
│ KEY        │ VALUE                                                                               │
├────────────┼─────────────────────────────────────────────────────────────────────────────────────┤
│ ID         │ ocid1.instance.oc1.iad.123                                                          │
├────────────┼─────────────────────────────────────────────────────────────────────────────────────┤
│ Name       │ my-instance                                                                         │
├────────────┼─────────────────────────────────────────────────────────────────────────────────────┤
│ State      │ RUNNING                                                                             │
├────────────┼─────────────────────────────────────────────────────────────────────────────────────┤
│ Region     │ us-ashburn-1                                                                        │
├────────────┼─────────────────────────────────────────────────────────────────────────────────────┤
│ Shape      │ VM.Standard.E5.Flex                                                                 │
╰────────────┴─────────────────────────────────────────────────────────────────────────────────────╯
```

In this example:
- The tenancy name is colored magenta
- The compartment name is colored cyan
- The instance name is colored blue
- The values are colored yellow

### Printing a Row-Based Table

To print a row-based table with headers:

```
headers := []string{"Name", "ID", "State", "Region"}
rows := [][]string{
    {"instance-1", "ocid1.instance.oc1.iad.123", "RUNNING", "us-ashburn-1"},
    {"instance-2", "ocid1.instance.oc1.iad.456", "STOPPED", "us-ashburn-1"},
    {"instance-3", "ocid1.instance.oc1.iad.789", "RUNNING", "us-phoenix-1"},
}
tablePrinter.PrintRowsTable(headers, rows)
```

This will produce output similar to:

```
┌───────────┬──────────────────────────┬─────────┬──────────────┐
│   NAME    │           ID             │  STATE  │    REGION    │
├───────────┼──────────────────────────┼─────────┼──────────────┤
│ instance-1│ ocid1.instance.oc1.iad.123│ RUNNING │ us-ashburn-1 │
│ instance-2│ ocid1.instance.oc1.iad.456│ STOPPED │ us-ashburn-1 │
│ instance-3│ ocid1.instance.oc1.iad.789│ RUNNING │ us-phoenix-1 │
└───────────┴──────────────────────────┴─────────┴──────────────┘
```

## Real-World Example

Here's a real-world example from the ocloud project that shows how to use the `TablePrinter` to display instance information:

```
func PrintInstancesTable(instances []Instance, appCtx *app.AppContext) {
    // Create a table printer with the tenancy name as the title
    tablePrinter := printer.NewTablePrinter(appCtx.TenancyName)

    // Convert instances to a format suitable for the printer
    if len(instances) == 0 {
        fmt.Println("No instances found.")
        return
    }

    // Print each instance as a key-value table with a title
    for _, instance := range instances {
        // Create a map with the instance data
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

        // Define the order of keys to match the example
        orderedKeys := []string{
            "ID",
            "AD",
            "FD",
            "Region",
            "Shape",
            "vCPUs",
            "Created",
            "Subnet ID",
            "Name",
            "Private IP",
            "Memory",
            "State",
        }

        // Print the table with ordered keys and colored title components
        tablePrinter.PrintKeyValueTableWithTitleOrdered(appCtx, instance.Name, instanceData, orderedKeys)
    }
}
```

## Styling Options

The `TablePrinter` uses the `StyleRounded` style from the go-pretty library for key-value tables and the `StyleLight` style for row-based tables. These styles provide a clean and readable output with rounded corners for key-value tables and light borders for row-based tables.

You can see all available styles in the [go-pretty documentation](https://github.com/jedib0t/go-pretty/blob/main/table/style.go).

## Color Options

The `PrintKeyValueTableWithTitleOrdered` method applies colors to different components of the table:

- Tenancy name: Magenta (`text.FgMagenta`)
- Compartment name: Cyan (`text.FgCyan`)
- Instance name: Blue (`text.FgBlue`)
- Values: Yellow (`text.FgYellow`)

These colors help to visually distinguish different components of the table and make the output more readable.

## Customization

While the `TablePrinter` provides sensible defaults, you can customize the output by modifying the table style, colors, and other properties. For example:

```
t := table.NewWriter()
t.SetOutputMirror(os.Stdout)
t.SetStyle(table.StyleDouble) // Use double borders
t.Style().Title.Align = text.AlignLeft // Align title to the left
t.Style().Options.DrawBorder = false // Don't draw the outer border
```

For more customization options, refer to the [go-pretty documentation](https://github.com/jedib0t/go-pretty).

## Conclusion

The `TablePrinter` package provides a simple and flexible way to display structured data in the terminal. It's particularly useful for displaying information about cloud resources, such as instances, in a readable and visually appealing format.
