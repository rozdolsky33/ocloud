# OCloud - Oracle Cloud Infrastructure CLI Tool

[![CI Build](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml/badge.svg)](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rozdolsky33/ocloud?sort=semver)
[![Downloads](https://img.shields.io/github/downloads/rozdolsky33/ocloud/total?style=flat&logo=cloudsmith&logoColor=white&label=downloads&labelColor=2f363d&color=brightgreen)](https://github.com/rozdolsky33/ocloud/releases)
[![Version](https://img.shields.io/badge/goversion-1.25.x-blue.svg)](https://golang.org)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/rozdolsky33/ocloud/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rozdolsky33/ocloud)](https://goreportcard.com/report/github.com/rozdolsky33/ocloud)
[![Go Coverage](https://github.com/rozdolsky33/ocloud/wiki/coverage.svg)](https://raw.githack.com/wiki/rozdolsky33/ocloud/coverage.html)

## Overview

OCloud is a powerful command-line interface (CLI) tool designed to simplify interactions with Oracle Cloud Infrastructure (OCI). It provides a streamlined experience for exploring and retrieving information about common OCI services with a focus on usability, performance, and automation.

Whether you're exploring instances, working with databases, or need to quickly find resources across your OCI environment, OCloud offers an intuitive and efficient interface that works seamlessly with your existing OCI configuration.

## Features

### Compute Resources
- **Instances**: List, search, and explore compute instances
- **Images**: Browse and search compute images
- **OKE Clusters**: Manage Kubernetes clusters and node pools

### Database Services
- **Autonomous Database**: List, search, and explore ADB instances
- **HeatWave MySQL**: Manage HeatWave database instances with detailed configuration info

### Networking
- **VCNs**: Virtual Cloud Networks with gateways, subnets, NSGs, route tables, and security lists
- **Subnets**: Network subnet management
- **Load Balancers**: Explore and search load balancer configurations

### Identity & Access
- **Compartments**: Navigate compartment hierarchy with tenancy-level scope support
- **Policies**: Explore IAM policies across compartments
- **Bastion Sessions**: Manage OCI Bastion sessions with interactive SSH key picker (TUI)

### Storage
- **Object Storage**: Browse and search buckets with interactive TUI

### Core Capabilities
- **Powerful Search**: Fuzzy, prefix, and substring matching using Bleve indexing
- **Interactive TUI**: Navigate resources with terminal user interface for select commands
- **JSON Output**: Consistent structured output across all commands
- **Pagination**: Unified pagination support (`--limit`, `--page`)
- **Authentication**: Interactive OCI Auth with automatic session refresh
- **Tenancy Mapping**: Friendly names for tenancies and compartments

## Installation

### Using Homebrew (macOS and Linux)

```bash
# Add the tap
brew tap rozdolsky33/ocloud https://github.com/rozdolsky33/ocloud

# Install ocloud
brew install ocloud
```

### Manual Installation

Download the latest binary from the [release page](https://github.com/rozdolsky33/ocloud/releases) and place it in your PATH.

**macOS/Linux:**
```bash
# Move the binary to a directory in your PATH
mv ~/Downloads/ocloud ~/.local/bin

# For macOS, clear quarantine (if applicable)
sudo xattr -d com.apple.quarantine ~/.local/bin/ocloud

# Make executable
chmod +x ~/.local/bin/ocloud
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/rozdolsky33/ocloud.git
cd ocloud

# Build and install
make build
make install
```

## Prerequisites

- **OCI CLI**: Installed and configured ([Quickstart guide](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm#Quickstart))
- **kubectl**: For OKE cluster interactions ([macOS](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/) | [Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/))
- **SSH Key Pair**: Required for OCI Bastion sessions (`~/.ssh`)

## Quick Start

```bash
# View configuration and available commands
ocloud

# Authenticate with OCI
ocloud config session authenticate

# Search for compute instances
ocloud compute instance search "prod"

# Get HeatWave databases with pagination
ocloud database heatwave get --limit 10 --page 1

# Search Autonomous Databases in JSON format
ocloud database autonomous search "test" --json

# Interactive VCN list (TUI)
ocloud network vcn list
```

## Configuration

OCloud can be configured via (precedence: highest to lowest):
1. Command-line flags
2. Environment variables
3. OCI configuration file (`~/.oci/config`)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OCI_CLI_PROFILE` | OCI configuration profile | `DEFAULT` |
| `OCI_TENANCY_NAME` | Tenancy name (looked up in tenancy map) | - |
| `OCI_COMPARTMENT` | Compartment name | - |
| `OCI_REGION` | OCI region | - |
| `OCI_TENANCY_MAP_PATH` | Path to tenancy mapping file | `~/.oci/.ocloud/tenancy-map.yaml` |

### Authentication

OCloud provides interactive authentication with automatic session refresh:

```bash
# Authenticate with OCI
ocloud config session authenticate

# Filter regions by prefix (e.g., us, eu, ap)
ocloud config session authenticate --filter us

# Filter by realm (e.g., OC1, OC2)
ocloud config session authenticate --realm OC1
```

The **OCI Auth Refresher** is a background service that automatically refreshes your OCI session before it expires, keeping your session active for long-running operations. Status is displayed in configuration output: `OCI_AUTH_AUTO_REFRESHER: ON [PID]`.

### Tenancy Mapping

Map tenancy names to OCIDs using `~/.oci/.ocloud/tenancy-map.yaml`:

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

View tenancy mapping:
```bash
ocloud config info map-file
ocloud config info map-file --realm OC1 --json
```

## Command Reference

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--compartment` | `-c` | OCI compartment name |
| `--tenancy-name` | | Tenancy name |
| `--debug` | `-d` | Enable debug logging |
| `--json` | `-j` | Output in JSON format |
| `--help` | `-h` | Display help |
| `--version` | `-v` | Print version |
| `--color` | | Enable colored output |
| `--log-level` | | Set verbosity (debug, info, warn, error) |

### Common Command Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | JSON output format |
| `--all` | `-A` | Include all related information |
| `--limit` | `-m` | Records per page (default: 20) |
| `--page` | `-p` | Page number (default: 1) |
| `--scope` | | `compartment` (default) or `tenancy` |
| `--tenancy-scope` | `-T` | Force tenancy-level scope |

### Scope Control

Some identity commands support compartment or tenancy scope:

```bash
# Compartments
ocloud identity compartment get                 # configured compartment (default)
ocloud identity compartment get --scope tenancy # whole tenancy
ocloud identity compartment get -T              # same as above

# Policies
ocloud identity policy get --scope compartment
ocloud identity policy search prod -T
```

### Network Resource Toggles

For VCN commands, include specific resources or use `--all`:

| Flag | Short | Description |
|------|-------|-------------|
| `--gateway` | `-G` | Internet/NAT gateways |
| `--subnet` | `-S` | Subnets |
| `--nsg` | `-N` | Network security groups |
| `--route-table` | `-R` | Route tables |
| `--security-list` | `-L` | Security lists |
| `--all` | `-A` | All of the above |

**Examples:**
```bash
# Get VCNs with all related resources
ocloud network vcn get --all
ocloud network vcn get -G -S -N -R -L

# Search VCNs with JSON output
ocloud network vcn search prod -A -j
```

## Command Categories

### Compute

```bash
# Instances
ocloud compute instance get
ocloud compute instance search "roster" --json
ocloud comp inst s "roster" -j

# Images
ocloud compute image get --limit 10
ocloud compute image search "Oracle-Linux"
ocloud comp img s "Oracle-Linux" -j

# OKE Clusters
ocloud compute oke get
ocloud compute oke search "orion" --json
```

### Database

```bash
# Autonomous Database
ocloud database autonomous get --limit 10 --page 1
ocloud database autonomous search "test" --json
ocloud db adb s "test" -j

# HeatWave MySQL
ocloud database heatwave get --all
ocloud database heatwave search "prod" --json
ocloud database heatwave list  # Interactive TUI
ocloud db hw s "8.4" -j
```

### Network

```bash
# VCNs
ocloud network vcn get --all
ocloud network vcn search "prod" -A -j
ocloud network vcn list  # Interactive TUI

# Load Balancers
ocloud network load-balancer get
ocloud network load-balancer search "prod" --all
ocloud net lb s "prod" -A -j

# Subnets
ocloud network subnet list
ocloud network subnet find "pub" --json
```

### Identity

```bash
# Compartments
ocloud identity compartment get
ocloud identity compartment get -T  # Tenancy scope
ocloud identity compartment search "sandbox" --json

# Policies
ocloud identity policy get
ocloud identity policy search "monitor" -T
```

### Storage

```bash
# Object Storage
ocloud storage object-storage get
ocloud storage object-storage search "prod" --json
ocloud storage object-storage list  # Interactive TUI
ocloud storage os s "prod" -j
```

## Development

### Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary to `bin/ocloud` |
| `make install` | Install to `$GOPATH/bin` |
| `make test` | Run all tests with coverage |
| `make fmt` | Format Go source files |
| `make vet` | Run go vet static analysis |
| `make lint` | Run golangci-lint |
| `make vuln` | Run govulncheck |
| `make clean` | Remove build artifacts |
| `make release` | Build for all platforms and create archives |

### Testing

Run the comprehensive test script that covers all command categories:

```bash
./test_ocloud.sh
```

The script tests:
- Root commands and global flags
- Configuration commands (info, map-file, session)
- Compute commands (instance, image, oke)
- Identity commands (compartment, policy)
- Network commands (subnet, vcn, load-balancer)
- Storage commands (object-storage)
- Database commands (autonomous, heatwave)

## Tips & Best Practices

- **Interactive TUI Commands**: Use `q`, `Esc`, or `Ctrl+C` to quit. Exiting without selection is graceful.
- **Search Patterns**: Searches support fuzzy matching, so typos and partial matches work.
- **JSON Output**: Use `--json` for scripting and automation.
- **Debug Logging**: Run with `--log-level debug` or `-d` for troubleshooting.
- **Colored Output**: Use `--color` for better readability in terminals.

## Error Handling

OCloud provides detailed error messages with configurable log levels:

- `--log-level debug`: Most verbose, shows all logs (equivalent to `-d`)
- `--log-level info`: Default, shows standard information
- `--log-level warn`: Warnings and errors only
- `--log-level error`: Errors only

Enable colored output with `--color` for better visibility.

## Contributing

We welcome contributions! See [CONTRIBUTING.md](https://github.com/rozdolsky33/ocloud/blob/main/CONTRIBUTING.md) for:
- Development setup and architecture overview
- Coding standards and testing requirements
- Data flow patterns across domain, mapping, services, and cmd layers

## Reporting Issues

Use our GitHub Issue Forms for faster triage:
- **Bug Report**: Include command(s), environment details, region/tenancy context, and debug logs
- **Feature Request**: Describe CLI UX impact and affected layers

Open a new [issue](https://github.com/rozdolsky33/ocloud/issues/new/choose)

**Tips for bug reports:**
- Run with `--log-level debug` and include relevant logs (redact secrets)
- Provide exact command(s), flags, OS/arch, and `ocloud version` output

For questions, use [Discussions](https://github.com/rozdolsky33/ocloud/discussions)

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.
