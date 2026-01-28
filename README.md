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
- **Instances**: List, search, and explore compute instances with interactive TUI
- **Images**: Browse and search compute images
- **OKE Clusters**: List, search, and explore Kubernetes clusters with node pool details

### Database Services
- **Autonomous Database**: List, search, and explore ADB instances with interactive TUI
- **HeatWave MySQL**: List, search, and explore HeatWave database instances with interactive TUI
- **OCI Cache Cluster**: List, search, and explore OCI Cache Clusters (Redis/Valkey) with interactive TUI

### Networking
- **VCNs**: Virtual Cloud Networks with gateways, subnets, NSGs, route tables, and security lists
- **Subnets**: Network subnet management
- **Load Balancers**: Explore and search load balancer configurations with health summaries

### Identity & Access
- **Compartments**: Navigate compartment hierarchy with tenancy-level scope support
- **Policies**: Explore IAM policies across compartments
- **Bastion**: Comprehensive bastion and session management
    - List and explore existing bastions
    - Create interactive bastion sessions with TUI-guided flows
    - Connect to Compute Instances (Managed SSH & Port Forwarding)
    - Connect to Databases (Autonomous DB & HeatWave via Port Forwarding)
    - Connect to OKE Clusters (Managed SSH to nodes & Port Forwarding to API server)
    - Connect to Load Balancers (Port Forwarding with TUI selection and health summaries)
    - Enhanced privileged port handling with sudo password validation
    - Automatic SSH tunnel management with background processes
    - Interactive SSH key pair selection
    - Automatic kubeconfig setup for OKE connections

### Storage
- **Object Storage**: Comprehensive management with interactive TUI
    - Browse and search buckets and objects
    - Interactive **upload** with file browser and multipart support
    - Interactive **download** with real-time progress tracking
    - Human-readable file sizes and visual progress bars

### Core Capabilities
- **Powerful Search**: Fuzzy, prefix, and substring matching using Bleve indexing
- **Interactive TUI**: Navigate resources with the terminal user interface for select commands
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

# Interactive bucket and object browsing
ocloud storage object-storage list

# Upload file to Object Storage (Interactive TUI)
ocloud storage object-storage upload

# Create bastion session with interactive TUI flow
ocloud identity bastion create

# List private Load Balancers and health status
ocloud network load-balancer list
```

## Configuration

Running `ocloud` without any arguments displays the configuration details and available commands.

Example output (values will vary by version, time, and your environment):

```
 ██████╗  ██████╗██╗      ██████╗ ██╗   ██╗██████╗
██╔═══██╗██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗
██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║
██║   ██║██║     ██║     ██║   ██║██║   ██║██║  ██║
╚██████╔╝╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝
 ╚═════╝  ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝

	      Version: v0.1.10

Configuration Details: Valid until <timestamp>
  OCI_CLI_PROFILE: DEFAULT
  OCI_TENANCY_NAME: cloudops
  OCI_COMPARTMENT_NAME: cnopslabsdev1
  OCI_AUTH_AUTO_REFRESHER: ON [<pid>]
  PORT_FORWARDING: ON [1521, 3306, 6443]
  OCI_TENANCY_MAP_PATH: /Users/<name>/.oci/.ocloud/tenancy-map.yaml

Interact with Oracle Cloud Infrastructure

Usage:
  ocloud [flags]
  ocloud [command]

Available Commands:
  compute     Explore OCI compute services
  config      Configure ocloud CLI and authentication
  database    Explore OCI Database services
  help        Help about any command
  identity    Explore OCI identity services
  network     Explore OCI networking services
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

## Bastion Session Management

OCloud provides comprehensive bastion session management with interactive TUI-guided flows for secure access to OCI resources.

### Creating Bastion Sessions

The `ocloud identity bastion create` command launches an interactive flow that guides you through:

1. **Session Type Selection**: Choose between Bastion management or creating a new session
2. **Bastion Selection**: Pick from your active bastions via TUI
3. **Target Type Selection**: Choose your connection target (Instance, Database, OKE, or Load Balancer)
4. **Session Type**: Select Managed SSH or Port Forwarding
5. **Resource Selection**: Interactive TUI to pick the specific resource
6. **SSH Key Selection**: Choose your SSH key pair from `~/.ssh`
7. **Connection Setup**: Automatic tunnel creation and configuration

### Supported Connection Types

#### Compute Instance Connections

**Managed SSH**: Direct interactive SSH session to compute instances
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → Instance → Managed SSH → Pick Instance → Select Keys
```

**Port Forwarding**: Create an SSH tunnel to instance ports (e.g., VNC, RDP, custom apps)
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → Instance → Port Forwarding → Pick Instance → Enter Port
# SSH tunnel runs in background, logs written to ~/.oci/.ocloud/logs/
```

#### Database Connections

**Autonomous Database**: Secure port forwarding to Autonomous DB private endpoints
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → Database → Autonomous → Pick ADB → Enter Port (default: 1521)
# Tunnel runs in background, connect to localhost:<port>
```

**HeatWave MySQL**: Secure port forwarding to HeatWave database instances
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → Database → HeatWave → Pick HeatWave DB → Enter Port (default: 3306)
# Tunnel runs in background, connect to localhost:<port>
```

#### OKE Cluster Connections

**Managed SSH to Node**: Direct SSH access to OKE worker nodes
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → OKE → Managed SSH → Pick Cluster → Pick Node → Select Keys
```

**Port Forwarding to API Server**: Access Kubernetes API server via bastion
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → OKE → Port Forwarding → Pick Cluster → Enter Port (default: 6443)
# Automatically offers to create/merge kubeconfig
# kubectl commands work via the tunnel to localhost:<port>
```

#### Load Balancer Connections

**Port Forwarding**: Secure access to private Load Balancers with TUI-guided selection and health summaries
```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → Load Balancer → Pick LB → Enter Port (default local: 8443, target: 443)
# Note: Supports privileged local ports (e.g., 443) with sudo password validation
```

### SSH Tunnel Management

- Tunnels run as **background processes** and persist after CLI exits
- **Logs** written to `~/.oci/sessions/<profile>/logs/ssh-tunnel-<port>-<date>.log`
- **State tracking** in `~/.oci/sessions/<profile>/tunnel/tunnel-<port>.json`
- Automatic **port availability** checking
- **Privileged port support**: Secure handling of ports < 1024 with sudo password validation
- **Connection verification** with a 30-second timeout

### Listing Existing Bastions

```bash
ocloud identity bastion get              # List all bastions in compartment
ocloud identity bastion get --json       # JSON output
```

## Command Categories

### Compute

```bash
# Instances
ocloud compute instance get
ocloud compute instance list  # Interactive TUI
ocloud compute instance search "roster" --json
ocloud comp inst s "roster" -j

# Images
ocloud compute image get --limit 10
ocloud compute image list  # Interactive TUI
ocloud compute image search "Oracle-Linux"
ocloud comp img s "Oracle-Linux" -j

# OKE Clusters
ocloud compute oke get
ocloud compute oke list  # Interactive TUI
ocloud compute oke search "orion" --json
```

### Database

```bash
# Autonomous Database
ocloud database autonomous get --limit 10 --page 1
ocloud database autonomous list  # Interactive TUI
ocloud database autonomous search "test" --json
ocloud db adb s "test" -j

# HeatWave MySQL
ocloud database heatwave get --all
ocloud database heatwave list  # Interactive TUI
ocloud database heatwave search "prod" --json
ocloud db hw s "8.4" -j

# OCI Cache Cluster (Redis/Valkey)
ocloud database cache-cluster get --all
ocloud database cache-cluster list  # Interactive TUI
ocloud database cache-cluster search "prod" --json
ocloud db cc s "VALKEY_7_2" -j
# Alternative aliases: cachecluster, cc
```

### Network

```bash
# VCNs
ocloud network vcn get --all
ocloud network vcn list  # Interactive TUI
ocloud network vcn search "prod" -A -j

# Load Balancers
ocloud network load-balancer get
ocloud network load-balancer list  # Interactive TUI
ocloud network load-balancer search "prod" --all
ocloud net lb s "prod" -A -j

# Subnets
ocloud network subnet list  # Interactive TUI
ocloud network subnet find "pub" --json
```

### Identity

```bash
# Compartments
ocloud identity compartment get
ocloud identity compartment list             # Interactive TUI
ocloud identity compartment get -T           # Tenancy scope
ocloud identity compartment search "sandbox" --json

# Policies
ocloud identity policy get
ocloud identity policy list                  # Interactive TUI
ocloud identity policy search "monitor" -T

# Bastion
ocloud identity bastion get                  # List existing bastions
ocloud identity bastion create               # Interactive bastion session creation
ocloud ident b create                        # Short alias
```

### Storage

```bash
# Object Storage
ocloud storage object-storage get      # List buckets
ocloud storage object-storage list     # Interactive TUI (browse buckets & objects)
ocloud storage object-storage search "prod" --json
ocloud storage object-storage upload   # Interactive TUI upload
ocloud storage object-storage download # Interactive TUI download
ocloud storage os s "prod" -j          # Search alias
```

## Practical Examples

### Secure Database Access via Bastion

Connect to a private Autonomous Database through a bastion session:

```bash
# 1. Authenticate with OCI
ocloud config session authenticate

# 2. Create bastion session with interactive flow
ocloud identity bastion create
# Select: Session → Choose Active Bastion → Database → Autonomous → Pick ADB
# The tunnel runs in the background on localhost:1521 (or custom port)

# 3. Connect using SQL client
sqlplus admin@localhost:1521/mydb_high
```

### OKE Cluster Management via Bastion

Access a private OKE cluster's API server through port forwarding:

```bash
# Create port forwarding session to OKE API server
ocloud identity bastion create
# Select: Session → Choose Bastion → OKE → Port Forwarding → Pick Cluster
# Accept kubeconfig creation/merge when prompted

# Tunnel runs in background, kubectl now works via localhost:6443
kubectl get nodes
kubectl get pods --all-namespaces
```

### SSH to Compute Instance via Bastion

Interactive SSH session to a private compute instance:

```bash
ocloud identity bastion create
# Select: Session → Choose Bastion → Instance → Managed SSH → Pick Instance
# Provides direct interactive SSH shell
```

### Search and Connect Workflow

Combine search with bastion for quick access:

```bash
# 1. Search for HeatWave databases
ocloud database heatwave search "prod" --json | jq '.[] | select(.lifecycleState=="ACTIVE")'

# 2. Create connection to selected database
ocloud identity bastion create
# Select: Session → Bastion → Database → HeatWave → Pick the prod DB
# Tunnel established to localhost:3306

# 3. Connect with MySQL client
mysql -h 127.0.0.1 -P 3306 -u admin -p
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
- Database commands (autonomous, heatwave, cache-cluster)

## Tips & Best Practices

- **Interactive TUI Commands**: Use `q`, `Esc`, or `Ctrl+C` to quit. Exiting without selection is graceful.
- **Search Patterns**: Searches support fuzzy matching, so typos and partial matches work.
- **JSON Output**: Use `--json` for scripting and automation.
- **Debug Logging**: Run with `--log-level debug` or `-d` for troubleshooting.
- **Colored Output**: Use `--color` for better readability in terminals.
- **Bastion Tunnels**: Background SSH tunnels persist after CLI exits. Check logs in `~/.oci/sessions/<profile>/logs/` if connection fails.
- **OKE Access**: Port forwarding to the OKE API server automatically offers kubeconfig setup for seamless kubectl access.
- **Database Connections**: Use bastion port forwarding for secure access to private Autonomous DB and HeatWave instances without exposing public endpoints.
- **Object Storage**: The `list` command allows interactive browsing of objects within buckets. Use `upload` and `download` for easy file transfers with progress feedback.

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


## Third-Party Attributions

OCloud uses third-party open-source software, including:

- **Oracle Cloud Infrastructure Go SDK** - Dual-licensed under UPL-1.0 or Apache-2.0
  Copyright (c) 2016, 2024, Oracle and/or its affiliates

For a complete list of third-party software and their licenses, see:
- [NOTICE](NOTICE) - Third-party software notices and attributions
- [third_party/](third_party/) - Full license texts for dependencies
