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

### Authentication

OCloud provides interactive authentication with OCI through the `config session` command:

```bash
# Authenticate with OCI
ocloud config session authenticate

# Filter regions by prefix (e.g., us, eu, ap)
ocloud config session authenticate --filter us

# Filter by realm (e.g., OC1, OC2)
ocloud config session authenticate --realm OC1
```

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

To view the tenancy mapping information:

```bash
# View tenancy mapping information
ocloud config info map-file

# View in JSON format
ocloud config info map-file --json

# Filter by realm
ocloud config info map-file --realm OC1
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `OCI_CLI_PROFILE` | OCI configuration profile (default: "DEFAULT") |
| `OCI_CLI_TENANCY` | Tenancy OCID |
| `OCI_TENANCY_NAME` | Tenancy name (looked up in tenancy map) |
| `OCI_COMPARTMENT` | Compartment name |
| `OCI_REGION` | OCI region |
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

#### Command Flags

| Flag      | Short | Description |
|-----------|-------|-------------|
| `--json`  | `-j`  | Output information in JSON format |
| `--all`   | `-A`  | Show all information |
| `--limit` | `-m`  | Maximum number of records per page (default: 20) |
| `--page`  | `-p`  | Page number to display (default: 1) |
| `--filter` | `-f` | Filter regions by prefix (e.g., us, eu, ap) |
| `--realm` | `-r` | Filter by realm (e.g., OC1, OC2) |

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
- Configuration commands (info, map-file, session)
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

You can also use the shorthand `-d` flag to enable debug logging.

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request