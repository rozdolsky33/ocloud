# OCloud - Oracle Cloud Infrastructure CLI Tool
[![CI Build](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml/badge.svg)](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rozdolsky33/oloud?sort=semver)
[![Go Coverage](https://github.com/rozdolsky33/ocloud/wiki/coverage.svg)](https://raw.githack.com/wiki/rozdolsky33/ocloud/coverage.html)
[![Version](https://img.shields.io/badge/goversion-1.24.x-blue.svg)](https://golang.org)
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

### Prerequisites

- Go 1.24 or later
- Oracle Cloud Infrastructure account
- OCI SDK configuration (typically in `~/.oci/config`)

### From Source

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
| `--compartment` | `-c` | OCI compartment name |

### Instance Command Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--find` | `-f` | Find resources by name pattern search |
| `--image-details` | `-i` | Show image details |
| `--list` | `-l` | List all resources |

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

# Show instance command help
ocloud instance --help

# Enable debug logging with colored output
ocloud --log-level debug --color instance list
```

### Working with Instances

```bash
# List all instances in a compartment
ocloud --compartment my-compartment instance list

# List all instances using the old flag syntax (for backward compatibility)
ocloud --compartment my-compartment instance -l

# Find instances by name pattern
ocloud --compartment my-compartment instance find "web-server"

# Find instances with image details
ocloud --compartment my-compartment instance find "web-server" --image-details
```

### Working with Different Tenancies

```bash
# Use a specific tenancy by OCID
ocloud --tenancy-id ocid1.tenancy.oc1..aaaaaaaa... instance list

# Use a tenancy by name (requires tenancy map)
ocloud --tenancy-name my-production-tenancy instance list

# Use environment variables
export OCI_TENANCY_NAME=my-production-tenancy
ocloud instance list
```

## Development

### Project Structure

The project follows a modern Go application structure:

- `cmd/`: Command-line interface implementation
  - `root.go`: Root command and global flags
  - `instance/`: Instance-related commands
    - `root.go`: Instance command and flags
    - `list.go`: List instances command
    - `find.go`: Find instances command
- `internal/`: Internal packages (not intended for external use)
  - `app/`: Application context and core functionality
  - `config/`: Configuration handling
  - `logger/`: Logging setup and utilities
  - `oci/`: OCI client factories
- `pkg/`: Public packages (can be imported by other projects)
  - `resources/`: Resource operations implementation

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

2. **go generate**: Used to keep flag constants, documentation, and code in sync. When flag definitions are changed in `pkg/flags/flags.go`, running `make generate` will automatically update the flag tables in this README.

To ensure your changes maintain code quality and consistency:

1. Run `make lint` before committing to check for code quality issues
2. Run `make generate` after modifying flag definitions to update documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
