# Printer Package

The `printer` package provides functionality to print data in tabular format using the [go-pretty](https://github.com/jedib0t/go-pretty) library. It offers a simple and flexible way to display structured data in the terminal with various formatting options.

## Overview

The `Printer` struct is the main component of this package, providing methods to print data in different formats:

- JSON output with proper indentation
- Key-value tables with ordered keys, titles, and colored values
- Multi-column tables with responsive column widths and styled headers

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

### Printing a Multi-Column Table with Responsive Widths

To print a multi-column table with headers and rows:

```
// Define table headers
headers := []string{"Name", "CIDR", "Public IP", "DNS Label", "Subnet Domain"}

// Create rows for the table
rows := [][]string{
    {"subnet-1", "10.0.0.0/24", "No", "subnet1", "subnet1.vcn.oraclevcn.com"},
    {"subnet-2", "10.0.1.0/24", "Yes", "subnet2", "subnet2.vcn.oraclevcn.com"},
    {"subnet-3", "10.0.2.0/24", "No", "subnet3", "subnet3.vcn.oraclevcn.com"},
}

// Print the table with a title
p.PrintTable("Subnet List", headers, rows)
```

This will produce output similar to:

```
╭──────────────────────────────────────────────────────────────────╮
│                         Subnet List                              │
├─────────┬─────────────┬───────────┬───────────┬──────────────────┤
│ NAME    │ CIDR        │ PUBLIC IP │ DNS LABEL │ SUBNET DOMAIN    │
├─────────┼─────────────┼───────────┼───────────┼──────────────────┤
│ subnet-1│ 10.0.0.0/24 │     No    │ subnet1   │ subnet1.vcn...   │
│ subnet-2│ 10.0.1.0/24 │    Yes    │ subnet2   │ subnet2.vcn...   │
│ subnet-3│ 10.0.2.0/24 │     No    │ subnet3   │ subnet3.vcn...   │
╰─────────┴─────────────┴───────────┴───────────┴──────────────────╯
```

In this example:
- Headers are displayed in bright yellow for better visibility
- Column widths automatically adjust based on terminal width
- Long text is truncated with ellipsis when necessary
- CIDR and IP columns are automatically centered for better readability
- The title is centered at the top of the table

## Real-World Examples

### Example 1: Using PrintKeyValues for Instance Details

Here's a real-world example from the ocloud project that shows how to use the `PrintKeyValues` method to display detailed instance information:

```
// PrintInstancesInfo displays detailed information for each instance.
func PrintInstancesInfo(instances []Instance, appCtx *app.ApplicationContext, useJSON bool) error {
    // Create a new printer that writes to the application's standard output.
    p := printer.New(appCtx.Stdout)

    // If JSON output is requested, use the printer to marshal the response.
    if useJSON {
        return marshalInstancesToJSON(p, instances, nil)
    }

    // Handle the case where no instances are found.
    if len(instances) == 0 {
        fmt.Fprintln(appCtx.Stdout, "No instances found.")
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

### Example 2: Using PrintTable for Subnet Listing

Here's another example that shows how to use the `PrintTable` method to display a list of subnets in a tabular format:

```
// PrintSubnetTable displays subnets in a formatted table or JSON format.
func PrintSubnetTable(subnets []Subnet, appCtx *app.ApplicationContext, pagination *util.PaginationInfo, useJSON bool, sortBy string) error {
    // Create a new printer that writes to the application's standard output.
    p := printer.New(appCtx.Stdout)

    // If JSON output is requested, use the printer to marshal the response.
    if useJSON {
        return util.MarshalDataToJSONResponse[Subnet](p, subnets, pagination)
    }

    // Handle the case where no subnets are found.
    if len(subnets) == 0 {
        fmt.Fprintln(appCtx.Stdout, "No Items found.")
        return nil
    }

    // Define table headers
    headers := []string{"Name", "CIDR", "Public IP", "DNS Label", "Subnet Domain"}

    // Create rows for the table
    rows := make([][]string, len(subnets))
    for i, subnet := range subnets {
        // Determine if public IP is allowed
        publicIPAllowed := "No"
        if !subnet.ProhibitPublicIPOnVnic {
            publicIPAllowed = "Yes"
        }

        // Create a row for this subnet
        rows[i] = []string{
            subnet.Name,
            subnet.CIDR,
            publicIPAllowed,
            subnet.DNSLabel,
            subnet.SubnetDomainName,
        }
    }

    // Create the colored title using components from the app context
    title := util.FormatColoredTitle(appCtx, "Subnets")

    // Print the table
    p.PrintTable(title, headers, rows)

    // Display pagination information if available
    if pagination != nil {
        fmt.Fprintf(appCtx.Stdout, "Page %d | Total: %d\n", pagination.CurrentPage, pagination.TotalCount)
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

The `PrintTable` method applies bright yellow color to the headers:

- Headers: Bright Yellow (`text.FgHiYellow`)

This coloring helps to visually distinguish the values from the keys and the headers from the data, making the output more readable.

In the real-world examples, additional colors are applied to the title components:

- Tenancy name: Magenta (`text.FgMagenta`)
- Compartment name: Cyan (`text.FgCyan`)
- Resource name (e.g., instance, subnet): Blue (`text.FgBlue`)

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

The `printer` package provides a simple and flexible way to display structured data in the terminal. It's particularly useful for displaying information about cloud resources, such as instances, compartments, and subnets, in a readable and visually appealing format. The package supports JSON output, key-value tables, and multi-column tables, making it versatile for different use cases.

The responsive design of the tables ensures that output looks good on different terminal sizes, and the automatic truncation of long text with ellipsis helps maintain readability. Special formatting for certain column types (like CIDR and IP addresses) further enhances the user experience.
