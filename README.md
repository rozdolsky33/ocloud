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
- Debug mode for troubleshooting

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

| Flag | Short | Description |
|------|-------|-------------|
| `--debug` | `-d` | Enable debug logging |
| `--tenancy-id` | `-t` | Tenancy OCID |
| `--compartment` | `-c` | Compartment name |
| `--list` | `-l` | List all resources |
| `--find` | `-f` | Find resources by name pattern search |
| `--create` | `-r` | Create a resource |
| `--type` | `-y` | Resource type |

Additional flags are available for specific operations. Use `ocloud --help` to see all available options.

## Usage Examples

### Basic Usage

```bash
# Show help
ocloud

# Enable debug mode
ocloud --debug

# List resources in a specific compartment
ocloud --compartment my-compartment --list
```

### Working with Different Tenancies

```bash
# Use a specific tenancy by OCID
ocloud --tenancy-id ocid1.tenancy.oc1..aaaaaaaa...

# Use a tenancy by name (requires tenancy map)
export OCI_TENANCY_NAME=my-production-tenancy
ocloud
```

## Development

### Project Structure

- `cmd/`: Command-line interface implementation
- `internal/`: Internal packages
  - `config/`: Configuration handling

### Development Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make run` | Build and run the CLI |
| `make test` | Run tests |
| `make fmt` | Format code |
| `make vet` | Run go vet |
| `make clean` | Clean build artifacts |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

[Contribution guidelines]
