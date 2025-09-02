# OCloud - Oracle Cloud Infrastructure CLI Tool
[![CI Build](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml/badge.svg)](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rozdolsky33/ocloud?sort=semver)
[![Downloads](https://img.shields.io/github/downloads/rozdolsky33/ocloud/total?style=flat&logo=cloudsmith&logoColor=white&label=downloads&labelColor=2f363d&color=brightgreen)](https://github.com/rozdolsky33/ocloud/releases)
[![Version](https://img.shields.io/badge/goversion-1.24.x-blue.svg)](https://golang.org)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/rozdolsky33/ocloud/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rozdolsky33/ocloud)](https://goreportcard.com/report/github.com/rozdolsky33/ocloud)
[![Go Coverage](https://github.com/rozdolsky33/ocloud/wiki/coverage.svg)](https://raw.githack.com/wiki/rozdolsky33/ocloud/coverage.html)

## Overview

OCloud is a powerful command-line interface (CLI) tool designed to simplify interactions with Oracle Cloud Infrastructure (OCI). It provides a streamlined experience for managing common OCI services—including compute, identity, networking, and database—with a focus on usability, performance, and automation.

Whether you're managing instances, working with images, or need to quickly find resources across your OCI environment, OCloud offers an intuitive and efficient interface that works seamlessly with your existing OCI configuration.

## Features

- Manage compute instances, images, and OKE resources
- Interactive configuration and OCI Auth Refresher to keep sessions alive
- Tenancy mapping for friendly tenancy and compartment names
- Bastion session management: start/attach/terminate OCI Bastion sessions with reachability checks and an interactive SSH key picker (TUI)

### What's New

- Added bastion session management capabilities and interactive SSH key selection to the Identity > Bastion commands

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

Download the latest binary from the [release page](https://github.com/rozdolsky33/ocloud/releases) and place it in your PATH.

#### macOS/Linux

```bash
# Move the binary to a directory in your PATH
mv ~/Downloads/ocloud ~/.local/bin

# For macOS, clear quarantine (if applicable)
sudo xattr -d com.apple.quarantine ~/.local/bin/ocloud

# Make the binary executable
chmod +x ~/.local/bin/ocloud
```


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

Running `ocloud` without any arguments displays the configuration details and available commands:

```
 ██████╗  ██████╗██╗      ██████╗ ██╗   ██╗██████╗
██╔═══██╗██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗
██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║
██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║
╚██████╔╝╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝
 ╚═════╝  ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝

	      Version: 0.0.31

Configuration Details: Valid until 2025-08-02 23:26:28
  OCI_CLI_PROFILE: DEFAULT
  OCI_TENANCY_NAME: cloudops
  OCI_COMPARTMENT_NAME: cnopslabsdev1
  OCI_AUTH_AUTO_REFRESHER: ON [44123]
  OCI_TENANCY_MAP_PATH: /Users/<name>/.oci/.ocloud/tenancy-map.yaml

Interact with Oracle Cloud Infrastructure

Usage:
  ocloud [flags]
  ocloud [command]

Available Commands:
  compute     Manage OCI compute services
  config      Manage ocloud CLI configurations file and authentication
  database    Manage OCI Database services
  help        Help about any command
  identity    Manage OCI identity services
  network     Manage OCI networking services
  version     Print the version information

Flags:
      --color                 Enable colored log messages.
  -c, --compartment string    OCI compartment name
  -d, --debug                 Enable debug logging
  -h, --help                  help for ocloud (shorthand: -h)
  -j, --json                  Output information in JSON format
      --log-level string      Set the log verbosity debug, (default "info")
  -t, --tenancy-id string     OCI tenancy OCID
      --tenancy-name string   Tenancy name
  -v, --version               Print the version number of ocloud CLI
```

OCloud can be configured in multiple ways, with the following precedence (highest to lowest):

1. Command-line flags
2. Environment variables
3. OCI configuration file

### OCI Configuration

OCloud uses the standard OCI configuration file located at `~/.oci/config`. You can specify a different profile using the `OCI_CLI_PROFILE` environment variable.

### Authentication

ocloud provides interactive authentication with OCI through the `config session` command:

```bash
# Authenticate with OCI
ocloud config session authenticate

# Filter regions by prefix (e.g., us, eu, ap)
ocloud config session authenticate --filter us

# Filter by realm (e.g., OC1, OC2)
ocloud config session authenticate --realm OC1
```

During authentication, you'll be prompted to set up the OCI Auth Refresher, which keeps your OCI session alive by automatically refreshing it before it expires. This is especially useful for long-running operations.

### OCI Auth Refresher

The OCI Auth Refresher is a background service that keeps your OCI session active by refreshing it shortly before it expires. When enabled, it runs as a background process and continuously monitors your session status.

Key features:
- Automatically refreshes your OCI session before it expires
- Runs in the background with minimal resource usage
- Supports multiple OCI profiles
- Status is displayed in the configuration output (`OCI_AUTH_AUTO_REFRESHER: ON [PID]`)

The refresher script is embedded in the ocloud binary and is automatically extracted to `~/.oci/.ocloud/scripts/` when needed.

### Tenancy Mapping

OCloud supports mapping tenancy names to OCIDs using a YAML file located at `~/.oci/.ocloud/tenancy-map.yaml`. The format is:

```yaml
- environment: Prod
  tenancy: cncloudps
  tenancy_id: ocid1.tenancy.oc1..aaaaaaaawdfste4i8fdsdsdkfasfds
  realm: OC1
  compartments:
    - sandbox
    - production 
  regions:
    - us-chicago-1
    - us-ashburn-1
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
| `OCI_AUTH_AUTO_REFRESHER` | Status of the OCI auth refresher |

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
- Identity commands (bastion, compartment, policy)
- Network commands (subnet)
- Database commands (autonomousdb)

The script tests various flags and abbreviations for each command, following a consistent pattern throughout.

To run the test script:

```bash
./test_ocloud.sh
```

## Error Handling

OCloud provides detailed error messages and supports multiple log verbosity levels. The valid values for --log-level are:

- `--log-level debug`: Most verbose; shows all logs, including detailed developer output. Equivalent to using `-d` / `--debug`.
- `--log-level info`: Default; shows standard, user‑facing information and errors.
- `--log-level warn`: Shows warnings and errors only.
- `--log-level error`: Shows errors only.

Tip: You can also enable colored log messages with `--color`.

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request