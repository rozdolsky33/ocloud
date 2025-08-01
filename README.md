# OCloud - Oracle Cloud Infrastructure CLI Tool
[![CI Build](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml/badge.svg)](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rozdolsky33/ocloud?sort=semver)
[![Version](https://img.shields.io/badge/goversion-1.21.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/rozdolsky33/ocloud/main/LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/rozdolsky33/ocloud)](https://goreportcard.com/report/github.com/rozdolsky33/ocloud)
[![Go Coverage](https://github.com/rozdolsky33/ocloud/wiki/coverage.svg)](https://raw.githack.com/wiki/rozdolsky33/ocloud/coverage.html)

## Overview

OCloud is a powerful command-line interface (CLI) tool designed to simplify interactions with Oracle Cloud Infrastructure (OCI). It provides a streamlined experience for managing compute resources with a focus on usability, performance, and automation capabilities.

Whether you're managing instances, working with images, or need to quickly find resources across your OCI environment, OCloud offers an intuitive and efficient interface that works seamlessly with your existing OCI configuration.

## Features

- **Compute Resource Management**
  - List and find compute instances with detailed information
  - View instance details including OS, image information, and resource specifications
  - Manage compute images with search capabilities
  - List and find Oracle Kubernetes Engine (OKE) clusters

- **Identity Management**
  - List and find compartments with detailed information
  - Manage identity policies with search capabilities

- **Database Management**
  - List and find Autonomous Databases with detailed information

- **Network Management**
  - List and find subnets with detailed information
  - Manage network resources with search capabilities

- **Configuration Management**
  - View and manage OCI configuration information
  - View tenancy mapping information
  - Authenticate with OCI and manage session tokens

- **Enhanced User Experience**
  - Colored output for improved readability
  - Configurable verbosity levels for debugging
  - JSON output support for automation and scripting
  - Pagination support for large result sets

- **Flexible Configuration**
  - Support for multiple tenancies and compartments
  - Configuration via environment variables, flags, or OCI config file
  - Tenancy mapping for simplified cross-environment management

- **Performance Optimizations**
  - Concurrent API calls for faster data retrieval (with rate limiting protection)
  - Efficient resource usage with pagination controls

## Installation

OCloud can be installed in several ways:

### Using Homebrew (macOS and Linux)

```bash
# Add the tap
brew tap rozdolsky33/ocloud https://github.com/rozdolsky33/ocloud

# Install ocloud
brew install ocloud
```

### Manual Installation

Download the latest binary from the [releases page](https://github.com/rozdolsky33/ocloud/releases) and place it in your PATH.

#### macOS/Linux

```bash
# Move the binary to a directory in your PATH
mv ~/Downloads/ocloud ~/.local/bin

# For macOS, clear quarantine (if applicable)
sudo xattr -d com.apple.quarantine ~/.local/bin/ocloud

# Make the binary executable
chmod +x ~/.local/bin/ocloud
```

#### Windows

1. Download the Windows binary from the releases page
2. Add the location to your PATH environment variable
3. Launch a new console session to apply the updated environment variable

### Build from Source

```bash
# Clone the repository
git clone https://github.com/rozdolsky33/ocloud.git
cd ocloud

# Build the binary
make build

# Install the binary to your GOPATH
make install
```

For detailed installation instructions, see the [Installation Guide](docs/installation.md).

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

### Environment Variables

| Variable | Description |
|----------|-------------|
| `OCI_CLI_PROFILE` | OCI configuration profile (default: "DEFAULT") |
| `OCI_CLI_TENANCY` | Tenancy OCID |
| `OCI_TENANCY_NAME` | Tenancy name (looked up in tenancy map) |
| `OCI_COMPARTMENT` | Compartment name |
| `OCI_CLI_REGION` | OCI region |
| `OCI_TENANCY_MAP_PATH` | Path to tenancy mapping file |

### Command-Line Flags

#### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--tenancy-id` | `-t` | OCI tenancy OCID |
| `--tenancy-name` |  | Tenancy name |
| `--log-level` |  | Set the log verbosity (debug, info, warn, error) |
| `--debug` | `-d` | Enable debug logging |
| `--color` |  | Enable colored output |
| `--compartment` | `-c` | OCI compartment name |
| `--disable-concurrency` | `-x` | Disable concurrency when fetching instance details |
| `--version` | `-v` | Print the version number |
| `--help` | `-h` | Display help information |

#### Instance Command Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output information in JSON format |
| `--image-details` | `-i` | Show image details including OS and tags |
| `--limit` | `-m` | Maximum number of records per page (default: 20) |
| `--page` | `-p` | Page number to display (default: 1) |

## Usage Examples

### Command Structure

OCloud uses a subcommand structure:

```
ocloud [global flags] <command> <subcommand> [command flags] [arguments]
```

### Basic Usage

```bash
# Show help
ocloud --help

# Show compute instance command help
ocloud compute instance --help

# Show version information
ocloud --version
```

### Working with Instances

```bash
# List all instances in a compartment
ocloud --compartment my-compartment compute instance list

# List all instances with pagination
ocloud --compartment my-compartment compute instance list --limit 10 --page 2

# Find instances by name pattern
ocloud --compartment my-compartment compute instance find "web-server"

# Find instances with image details
ocloud --compartment my-compartment compute instance find "web-server" --image-details

# Output instance information in JSON format
ocloud --compartment my-compartment compute instance list --json

# Find instances and output results in JSON format with image details
ocloud --compartment my-compartment compute instance find "web-server" --json --image-details
```

### Working with Images

```bash
# List all image in a compartment
ocloud --compartment my-compartment compute image list

# Find image by name pattern
ocloud --compartment my-compartment compute image find "Oracle-Linux"

# Output image information in JSON format
ocloud --compartment my-compartment compute image list --json
```

### Working with OKE Clusters

```bash
# List all OKE clusters in a compartment
ocloud --compartment my-compartment compute oke list

# Find OKE clusters by name pattern
ocloud --compartment my-compartment compute oke find "my-cluster"

# Output OKE cluster information in JSON format
ocloud --compartment my-compartment compute oke list --json
```

### Working with Compartments

```bash
# List all compartments
ocloud identity compartment list

# Find compartments by name pattern
ocloud identity compartment find "my-compartment"

# Output compartment information in JSON format
ocloud identity compartment list --json
```

### Working with Policies

```bash
# List all policies
ocloud identity policy list

# Find policies by name pattern
ocloud identity policy find "my-policy"

# Output policy information in JSON format
ocloud identity policy list --json
```

### Working with Autonomous Databases

```bash
# List all Autonomous Databases in a compartment
ocloud --compartment my-compartment database autonomousdb list

# Find Autonomous Databases by name pattern
ocloud --compartment my-compartment database autonomousdb find "my-database"

# Output Autonomous Database information in JSON format
ocloud --compartment my-compartment database autonomousdb list --json
```

### Working with Network Subnets

```bash
# List all subnets in a compartment
ocloud --compartment my-compartment network subnet list

# Find subnets by name pattern
ocloud --compartment my-compartment network subnet find "public"

# Output subnet information in JSON format
ocloud --compartment my-compartment network subnet list --json

# Use abbreviated commands
ocloud --compartment my-compartment net sub list
```

### Working with Configuration

```bash
# View tenancy mapping information
ocloud config info map-file

# View tenancy mapping information in JSON format
ocloud config info map-file --json

# Filter tenancy mapping by realm
ocloud config info map-file --realm OC1

# Authenticate with OCI (if implemented)
ocloud config auth
```

### Working with Different Tenancies

```bash
# Use a specific tenancy by OCID
ocloud --tenancy-id ocid1.tenancy.oc1..aaaaaaaa... compute instance list

# Use a tenancy by name (requires tenancy map)
ocloud --tenancy-name my-production-tenancy compute instance list

# Use environment variables
export OCI_TENANCY_NAME=my-production-tenancy
ocloud compute instance list
```

### Development Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make run` | Build and run the CLI |
| `make install` | Install the binary to $GOPATH/bin |
| `make test` | Run tests |
| `make fmt` | Format code |
| `make fmt-check` | Check if code is formatted correctly |
| `make vet` | Run go vet |
| `make lint` | Run golangci-lint |
| `make clean` | Clean build artifacts |
| `make release` | Build binaries for all supported platforms and create zip archives |
| `make compile` | Compile binaries for all supported platforms |
| `make zip` | Create zip archives for all binaries |
| `make check-env` | Check if required tools are installed |
| `make help` | Display help message with available commands |

### Testing

The project includes a comprehensive test script `test_ocloud.sh` that tests all major command categories and their subcommands:

- Root commands and global flags
- Configuration commands (info, map-file)
- Compute commands (instance, image, oke)
- Identity commands (compartment, policy)
- Network commands (subnet)
- Database commands (autonomousdb)

The script tests various flags and abbreviations for each command, following a consistent pattern throughout.

To run the test script:

```bash
./test_ocloud.sh
```

## Error Handling

OCloud provides detailed error messages and supports different verbosity levels:

- `--log-level info`: Shows standard information (default)
- `--log-level debug`: Shows detailed debugging information
- `--log-level warn`: Shows only warnings and errors
- `--log-level error`: Shows only errors

You can also use the shorthand `-d` flag to enable debug logging:

```bash
ocloud -d compute instance list
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
