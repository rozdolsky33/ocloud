# OCloud - Oracle Cloud Infrastructure CLI Tool
[![CI Build](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml/badge.svg)](https://github.com/rozdolsky33/ocloud/actions/workflows/build.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rozdolsky33/ocloud?sort=semver)
[![Downloads](https://img.shields.io/github/downloads/rozdolsky33/ocloud/total?style=flat&logo=cloudsmith&logoColor=white&label=downloads&labelColor=2f363d&color=brightgreen)](https://github.com/rozdolsky33/ocloud/releases)
[![Version](https://img.shields.io/badge/goversion-1.25.x-blue.svg)](https://golang.org)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/rozdolsky33/ocloud/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rozdolsky33/ocloud)](https://goreportcard.com/report/github.com/rozdolsky33/ocloud)
[![Go Coverage](https://github.com/rozdolsky33/ocloud/wiki/coverage.svg)](https://raw.githack.com/wiki/rozdolsky33/ocloud/coverage.html)

## Overview

OCloud is a powerful command-line interface (CLI) tool designed to simplify interactions with Oracle Cloud Infrastructure (OCI). It provides a streamlined experience for managing common OCI services—including compute, identity, networking, and database—with a focus on usability, performance, and automation.

Whether you're managing instances, working with images, or need to quickly find resources across your OCI environment, OCloud offers an intuitive and efficient interface that works seamlessly with your existing OCI configuration.

## Features

- Manage compute resources: Instances, Images, and OKE (Kubernetes) clusters and node pools
- Manage networking resources: Virtual Cloud Networks (VCNs), Subnets, Load Balancers, and related components
- Manage storage resources: Object Storage Buckets (list via TUI and paginated get)
- Powerful find/search commands using Bleve for fuzzy, prefix, and substring matching where applicable
- Interactive configuration and an OCI Auth Refresher to keep sessions alive
- Tenancy mapping for friendly tenancy and compartment names
- Bastion session management: start/attach/terminate OCI Bastion sessions with reachability checks and an interactive SSH key picker (TUI)
- Consistent JSON output, unified pagination across services, and short/long flag aliases
- Built-in security in CI: govulncheck vulnerability scanning via `make vuln` and GitHub Actions

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


## Prerequisites

Before using OCloud, ensure the following tools and setup are in place:
- OCI CLI installed and configured. Follow Oracle's Quickstart guide: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm#Quickstart
- kubectl installed (for interacting with OKE clusters). Installation instructions for Linux: https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
- SSH key pair available in your ~/.ssh directory. This is required for OCI Bastion managed sessions and for SSH port forwarding features.

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

	      Version: v0.1.0

Configuration Details: Valid until <timestamp>
  OCI_CLI_PROFILE: DEFAULT
  OCI_TENANCY_NAME: cloudops
  OCI_COMPARTMENT_NAME: cnopslabsdev1
  OCI_AUTH_AUTO_REFRESHER: ON [<pid>]
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

ocloud provides interactive authentication with OCI through the `config session` command. You can also control which browser is used during the login flow; see docs/auth-browser-selection.md for details:

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
| `--log-level` |  | Set the log verbosity (e.g., info, debug) |
| `--debug` | `-d` | Enable debug logging |
| `--color` |  | Enable colored output |
| `--compartment` | `-c` | OCI compartment name |
| `--version` | `-v` | Print the version number |
| `--help` | `-h` | Display help information |

#### Command Flags

| Flag            | Short | Description |
|-----------------|-------|-------------|
| `--json`        | `-j`  | Output information in JSON format |
| `--all`         | `-A`  | Include all related info sections where applicable (e.g., for VCNs: gateways, subnets, NSGs, route tables, security lists) |
| `--limit`       | `-m`  | Maximum number of records per page (default: 20) |
| `--page`        | `-p`  | Page number to display (default: 1) |
| `--scope`       |       | Listing/search scope for applicable commands: `compartment` (default) or `tenancy` |
| `--tenancy-scope` | `-T`  | Shortcut to force tenancy-level scope; overrides `--scope` |
| `--filter`      | `-f`  | Filter regions by prefix (e.g., us, eu, ap) |
| `--realm`       | `-r`  | Filter by realm (e.g., OC1, OC2) |

#### Network resource toggles (used by networking commands)

| Flag                | Short | Description                                 |
|---------------------|-------|---------------------------------------------|
| `--gateway`         | `-G`  | Include/display internet/NAT gateways       |
| `--subnet`          | `-S`  | Include/display subnets                     |
| `--nsg`             | `-N`  | Include/display network security groups     |
| `--route-table`     | `-R`  | Include/display route tables                |
| `--security-list`   | `-L`  | Include/display security lists              |

### Scope Control (Identity commands)

Some identity subcommands support scoping their operations to either the configured parent compartment or the whole tenancy.

- --scope: Choose where to operate: "compartment" (default) or "tenancy".
- -T/--tenancy-scope: Shortcut to force tenancy-level scope; it overrides --scope.

Examples:

```bash
# Compartments
ocloud identity compartment get                 # children of configured compartment (default)
ocloud identity compartment get --scope tenancy # whole tenancy (includes subtree)
ocloud identity compartment get -T              # same as above

# Policies
ocloud identity policy list --scope compartment # explicit compartment-level listing
ocloud identity policy search prod -T            # tenancy-level search
```

### Networking: VCN commands

The network VCN group provides commands to get and find Virtual Cloud Networks in the configured compartment. You can include related networking resources using the network toggles shown above or the --all (-A) flag to include everything at once.

Examples:

- Get VCNs with pagination
  - ocloud network vcn get
  - ocloud network vcn get --limit 10 --page 2
  - ocloud network vcn get -m 5 -p 3 --all
  - ocloud network vcn get -m 5 -p 3 -A -j

- Search VCNs by pattern
  - ocloud network vcn search prod
  - ocloud network vcn search prod --all
  - ocloud network vcn search prod -A -j

Interactive list (TUI):
- ocloud network vcn list
  Note: This command is interactive and not suitable for non-interactive scripts. If you quit without selecting an item, it exits without error.

### Networking: Load Balancer commands

Manage and explore Load Balancers in the configured compartment. You can:
- Get paginated lists with optional extra columns using --all (-A)
- Search using fuzzy, prefix, token, and substring matching across multiple fields
- Launch an interactive list (TUI) to search and select a Load Balancer

Examples:

- Get Load Balancers with pagination
  - ocloud network load-balancer get
  - ocloud network load-balancer get --limit 10 --page 2
  - ocloud network load-balancer get --all
  - ocloud net lb get -A -j

- Search Load Balancers by pattern
  - ocloud network load-balancer search prod
  - ocloud network load-balancer search prod --json
  - ocloud network load-balancer search prod --all
  - ocloud net lb s prod -A -j

Interactive list (TUI):
- ocloud network load-balancer list
  Note: This command is interactive and not suitable for non-interactive scripts. If you quit without selecting an item, it exits without error.

### Storage: Object Storage commands

Manage and explore Object Storage Buckets in the configured compartment. You can:
- Get paginated lists with optional extended details using --all (-A)
- Launch an interactive list (TUI) to search and select a Bucket

Examples:

- Get Buckets with pagination
  - ocloud storage object-storage get
  - ocloud storage object-storage get --limit 10 --page 2
  - ocloud storage object-storage get --all
  - ocloud storage os get -A -j

- Search Buckets by pattern
  - ocloud storage object-storage search prod
  - ocloud storage object-storage search prod --json
  - ocloud storage os s prod -j

Interactive list (TUI):
- ocloud storage object-storage list
  Note: This command is interactive and not suitable for non-interactive scripts. If you quit without selecting an item, it exits without error.

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
| `make vuln` | Run govulncheck vulnerability scan |
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
- Network commands (subnet, vcn, loadbalancer)
- Storage commands (object-storage)
- Database commands (autonomousdb)

The script tests various flags and abbreviations for each command, following a consistent pattern throughout.

To run the test script:

```bash
./test_ocloud.sh
```

## Tips

- Interactive TUI lists (e.g., network vcn list, database autonomous list, compute image list) support quitting with q/Esc/Ctrl+C. If you exit without selecting an item, the command will exit gracefully without an error.

## Error Handling

OCloud provides detailed error messages and supports multiple log verbosity levels. The valid values for --log-level are:

- `--log-level debug`: Most verbose; shows all logs, including detailed developer output. Equivalent to using `-d` / `--debug`.
- `--log-level info`: Default; shows standard, user‑facing information and errors.
- `--log-level warn`: Shows warnings and errors only.
- `--log-level error`: Shows errors only.

Tip: You can also enable colored log messages with `--color`.

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.