# OCloud - Oracle Cloud Infrastructure CLI Tool
[![CI Build](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml/badge.svg)](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rozdolsky33/ocloud?sort=semver)
[![Version](https://img.shields.io/badge/goversion-1.23.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/rozdolsky33/ocloud/main/LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/rozdolsky33/ocloud)](https://goreportcard.com/report/github.com/rozdolsky33/ocloud)

OCloud is a command-line interface (CLI) tool designed to simplify interactions with Oracle Cloud Infrastructure (OCI). It provides a streamlined experience for common OCI operations with a focus on usability and automation.

## Features

- Interact with Oracle Cloud Infrastructure resources
- Support for multiple tenancies and compartments
- Configuration via environment variables, flags, or OCI config file
- Colored logging with configurable verbosity levels
- Modular architecture for easy extension

## Installation

OCloud can be installed in several ways:

- Using Homebrew (macOS and Linux)
- Downloading pre-built binaries
- Building from source

For detailed installation instructions for all platforms, see the [Installation Guide](docs/installation.md).

### Prerequisites

- Go 1.21 or later
- Oracle Cloud Infrastructure account
- OCI SDK configuration (typically in `~/.oci/config`)

### Quick Start: Build From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/rozdolsky33/ocloud.git
   cd ocloud
   ```

2. Build the binary:
   ```bash
   make build
   ```

3. Install the binary to your GOPATH:
   ```bash
   make install
   ```

## Configuration

OCloud can be configured in multiple ways, with the following precedence (highest to lowest):

1. Command-line flags
2. Environment variables
3. OCI configuration file

### OCI Configuration

OCloud uses the standard OCI configuration file located at `~/.oci/config`. You can specify a different profile using the `OCI_CLI_PROFILE` environment variable.

### Tenancy Mapping

OCloud supports mapping tenancy names to OCIDs using a YAML file located at `~/.oci/tenancy-map.yaml`. The format is:

```yaml
- environment: "prod"
  tenancy: "my-production-tenancy"
  tenancy_id: "ocid1.tenancy.oc1..aaaaaaaa..."
  realm: "oc1"
  compartments: "compartment1,compartment2"
  regions: "us-ashburn-1,us-phoenix-1"
```

You can override the tenancy map path using the `OCI_TENANCY_MAP_PATH` environment variable.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `OCI_CLI_PROFILE` | OCI configuration profile (default: "DEFAULT") |
| `OCI_CLI_TENANCY` | Tenancy OCID |
| `OCI_TENANCY_NAME` | Tenancy name (looked up in tenancy map) |
| `OCI_COMPARTMENT` | Compartment name |
| `OCI_CLI_REGION` | OCI region |
| `OCI_TENANCY_MAP_PATH` | Path to tenancy mapping file |

## Command-Line Flags

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--tenancy-id` | `-t` | OCI tenancy OCID |
| `--tenancy-name` |  | Tenancy name |
| `--log-level` |  | Set the log verbosity. Supported values are: debug, info, warn, and error. |
| `--debug` | `-d` | Enable debug logging |
| `--color` |  | Enable colored output |
| `--compartment` | `-c` | OCI compartment name |
| `--disable-concurrency` | `-x` | Enable concurrency when fetching instance details (use -x to enable concurrency, which is disabled by default to avoid rate limiting) |
| `--version` | `-v` | Print the version number of ocloud CLI |

### Instance Command Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--find` | `-f` | Find resources by name pattern search |
| `--image-details` | `-i` | Show image details |
| `--json` | `-j` | Output information in JSON format |
| `--list` | `-l` | List all resources |
| `--limit` | `-m` | Maximum number of records to display per page (default: 20) |
| `--page` | `-p` | Page number to display (default: 1) |

Additional flags are available for specific operations. Use `ocloud --help` to see all available options.

## Usage Examples

### Command Structure

OCloud uses a subcommand structure:

```
ocloud [global flags] <command> [command flags] [arguments]
```

### Basic Usage

```bash
# Show help
ocloud --help

# Show compute instance command help
ocloud compute instance --help

# Show version information
ocloud --version
ocloud -v

# Enable debug logging with colored output
ocloud compute instance --list -d --color 

# Alternative way to enable debug logging
ocloud -d --color compute instance --list
```

### Working with Instances

```bash
# List all instances in a compartment
ocloud --compartment my-compartment compute instance --list

# List all instances using the shorthand flag syntax
ocloud --compartment my-compartment compute instance -l

# Find instances by name pattern
ocloud --compartment my-compartment compute instance --find "web-server"

# Find instances with image details
ocloud --compartment my-compartment compute instance --find "web-server" --image-details

# Output instance information in JSON format
ocloud --compartment my-compartment compute instance --list --json

# Output instance information in JSON format using shorthand flags
ocloud --compartment my-compartment compute instance -l -j

# Find instances and output results in JSON format
ocloud --compartment my-compartment compute instance --find "web-server" --json

# Enable concurrency for faster processing (may hit rate limits)
ocloud --compartment my-compartment -x compute instance --list

# Enable concurrency using the full flag name
ocloud --compartment my-compartment --disable-concurrency compute instance --find "web-server"

# List instances with pagination (default: 20 records per page)
ocloud --compartment my-compartment compute instance --list

# List instances with custom page size
ocloud --compartment my-compartment compute instance --list --limit 10

# Navigate to a specific page
ocloud --compartment my-compartment compute instance --list --page 2

# Combine pagination parameters
ocloud --compartment my-compartment compute instance --list --limit 5 --page 3

# Combine JSON output with pagination
ocloud --compartment my-compartment compute instance --list --limit 5 --page 2 --json
```

### Working with Different Tenancies

```bash
# Use a specific tenancy by OCID
ocloud --tenancy-id ocid1.tenancy.oc1..aaaaaaaa... compute instance --list

# Use a tenancy by name (requires tenancy map)
ocloud --tenancy-name my-production-tenancy compute instance --list

# Use environment variables
export OCI_TENANCY_NAME=my-production-tenancy
ocloud compute instance --list
```

## Development

### Project Structure

The project follows a modern Go application structure:

- `buildinfo/`: Version information
- `cmd/`: Command-line interface implementation
  - `root.go`: Root command and global flags
  - `compute/`: Compute-related commands
    - `root.go`: Compute command and flags
    - `instance/`: Instance-related commands
      - `root.go`: Instance command and flags
  - `version/`: Version command implementation
- `internal/`: Internal packages (not intended for external use)
  - `app/`: Application context and core functionality
  - `config/`: Configuration handling
    - `flags/`: CLI flag definitions and handling
    - `generator/`: Configuration generator utilities
  - `logger/`: Logging setup and utilities
  - `oci/`: OCI client factories
  - `printer/`: Output formatting utilities
  - `services/`: Service implementations
    - `compute/`: Compute resource operations
    - `database/`: Database resource operations
    - `identity/`: Identity resource operations
    - `network/`: Network resource operations

## Error Handling

OCloud provides detailed error messages and supports different verbosity levels:

- `--log-level info`: Shows standard information (default)
- `--log-level debug`: Shows detailed debugging information
- `--log-level warn`: Shows only warnings and errors
- `--log-level error`: Shows only errors

### Development Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make run` | Build and run the CLI |
| `make test` | Run tests |
| `make fmt` | Format code |
| `make vet` | Run go vet |
| `make lint` | Run golangci-lint |
| `make generate` | Run go generate to update generated code |
| `make clean` | Clean build artifacts |

### Code Quality and Drift Guard

OCloud uses several tools to maintain code quality and prevent drift between flag definitions, documentation, and code:

1. **golangci-lint**: A fast, parallel runner for Go linters. It helps catch issues like unused variables, formatting errors, and more. The configuration is in `.golangci.yml`.

2. **go generate**: Used to keep flag constants, documentation, and code in sync. When flag definitions are changed in `internal/config/flags/`, running `make generate` will automatically update the flag tables in this README.

To ensure your changes maintain code quality and consistency:

1. Run `make lint` before committing to check for code quality issues
2. Run `make generate` after modifying flag definitions to update documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
