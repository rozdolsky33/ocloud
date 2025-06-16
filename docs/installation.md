# Installation Guide

## Install `ocloud` using Homebrew

`ocloud` can be installed using [Homebrew](https://brew.sh) for macOS and Linux.

1. Add the `rozdolsky33/ocloud` tap:
   ```bash
   brew tap rozdolsky33/ocloud https://github.com/rozdolsky33/ocloud
   ```

2. Install `ocloud`:
   ```bash
   brew install ocloud
   ```

3. Verify the installation:
   ```bash
   ocloud -h
   ```

## Manual Installation

You can download the latest binary from the [ocloud releases](https://github.com/rozdolsky33/ocloud/releases) page.

### Place the Binary in Your `PATH`

#### macOS/Linux

After downloading the binary, move it to a directory included in your `PATH` (e.g., `/usr/local/bin` or any custom location).

**Example: Adding Binary to a Custom Path**

1. Check your `PATH`:
   ```bash
   echo $PATH
   ```
   Example output:
   ```
   /usr/local/bin:/Users/YOUR_USER/.local/bin
   ```

2. Move the binary to a directory in your `PATH` (e.g., `~/.local/bin`):
   ```bash
   mv ~/Downloads/ocloud ~/.local/bin
   ```

3. For macOS, clear quarantine (if applicable):
   ```bash
   sudo xattr -d com.apple.quarantine ~/.local/bin/ocloud
   ```

4. Make the binary executable:
   ```bash
   chmod +x ~/.local/bin/ocloud
   ```

5. Verify the installation:
   ```bash
   ocloud -h
   ```

#### Windows Setup

To use `ocloud` on Windows, you need to add its location to the `PATH` environment variable.

Steps:

1. Open **Control Panel** → **System** → **System Settings** → **Environment Variables**.
2. Scroll down in the **System Variables** section and locate the `PATH` variable.
3. Click **Edit** and add the location of your `ocloud` binary to the `PATH` variable. For example, `c:\ocloud`.

   *Note: When adding a new location, ensure that a semicolon (`;`) is included as a delimiter if appending to existing entries. Example: `c:\path;c:\ocloud`.*

4. Launch a new console session to apply the updated environment variable.

## Build from Source

If you prefer to build from source, follow these steps:

1. Prerequisites:
   - Go 1.20 or later
   - Git

2. Clone the repository:
   ```bash
   git clone https://github.com/rozdolsky33/ocloud.git
   cd ocloud
   ```

3. Build the binary:
   ```bash
   make build
   ```
   This will create the binary in the `bin` directory.

4. Install the binary to your GOPATH:
   ```bash
   make install
   ```

## Verify Installation

Once your installation is complete, verify it by running:

```bash
ocloud --version
```

or

```bash
ocloud -v
```

This should display version information similar to:
```
Version:    v0.1.0
Commit:     abcdef123456
Built:      2023-06-15 12:34:56
```

You can also check the help information:

```bash
ocloud --help
```

If the command prints the `ocloud` help information, the setup is complete.

## Configuration

After installation, you'll need to configure `ocloud` to work with your Oracle Cloud Infrastructure account. See the [Configuration](../README.md#configuration) section in the main README for details.