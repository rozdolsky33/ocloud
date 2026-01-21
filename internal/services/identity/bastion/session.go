package bastion

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

// Defaults used for session wait and ttl
var (
	waitPollInterval = 3 * time.Second
	defaultTTL       = 10800 // seconds (3 hours)
)

// TunnelInfo stores information about an active SSH tunnel
type TunnelInfo struct {
	PID       int       `json:"pid"`
	LocalPort int       `json:"local_port"`
	TargetIP  string    `json:"target_ip"`
	StartedAt time.Time `json:"started_at"`
	LogFile   string    `json:"log_file"`
}

// getTunnelsDir returns the directory where tunnel state files are stored
func getTunnelsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	profile := os.Getenv(flags.EnvKeyProfile)
	tunnelsDir := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCISessionsDirName, profile, "tunnels")
	return tunnelsDir, nil
}

// SaveTunnelState saves tunnel information to a state file
func SaveTunnelState(tunnel TunnelInfo) error {
	tunnelsDir, err := getTunnelsDir()
	if err != nil {
		return fmt.Errorf("get tunnels dir: %w", err)
	}

	if err := os.MkdirAll(tunnelsDir, 0o755); err != nil {
		return fmt.Errorf("create tunnels dir: %w", err)
	}

	stateFile := filepath.Join(tunnelsDir, fmt.Sprintf("tunnel-%d.json", tunnel.LocalPort))
	data, err := json.MarshalIndent(tunnel, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tunnel state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0o644); err != nil {
		return fmt.Errorf("write tunnel state: %w", err)
	}

	return nil
}

// GetActiveTunnels returns a list of currently active SSH tunnels
// It checks both state files and running processes for backward compatibility
func GetActiveTunnels() ([]TunnelInfo, error) {
	tunnelsMap := make(map[int]TunnelInfo)

	// 1. First, load tunnels from state files
	tunnelsDir, err := getTunnelsDir()
	if err == nil {
		if _, err := os.Stat(tunnelsDir); err == nil {
			entries, err := os.ReadDir(tunnelsDir)
			if err == nil {
				for _, entry := range entries {
					if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
						continue
					}

					stateFile := filepath.Join(tunnelsDir, entry.Name())
					data, err := os.ReadFile(stateFile)
					if err != nil {
						continue
					}

					var tunnel TunnelInfo
					if err := json.Unmarshal(data, &tunnel); err != nil {
						_ = os.Remove(stateFile)
						continue
					}

					if isProcessRunning(tunnel.PID) {
						tunnelsMap[tunnel.LocalPort] = tunnel
					} else {
						_ = os.Remove(stateFile)
					}
				}
			}
		}
	}

	// 2.Detect running SSH tunnels from a process list
	// This catches tunnels created before state file tracking was implemented
	cmd := exec.Command("pgrep", "-fl", "ssh")
	out, err := cmd.CombinedOutput()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			// Look for SSH port forwarding: "ssh ... -N ... -L port:target:port ..."
			if !strings.Contains(line, "-N") || !strings.Contains(line, "-L") {
				continue
			}
			// Parse PID and local port
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			pid, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}

			localPort := extractLocalPortFromSSHCommand(line)
			if localPort == 0 {
				continue
			}

			if _, exists := tunnelsMap[localPort]; !exists {
				tunnelsMap[localPort] = TunnelInfo{
					PID:       pid,
					LocalPort: localPort,
					TargetIP:  "unknown",
					StartedAt: time.Time{},
					LogFile:   "",
				}
			}
		}
	}

	// Convert map to slice
	var activeTunnels []TunnelInfo
	for _, tunnel := range tunnelsMap {
		activeTunnels = append(activeTunnels, tunnel)
	}

	return activeTunnels, nil
}

// extractLocalPortFromSSHCommand parses an SSH command line to extract the local port from -L flag
// Example: "-L 3306:10.0.0.156:3306" returns 3306
func extractLocalPortFromSSHCommand(cmdLine string) int {
	// Find the -L flag and extract the port
	parts := strings.Fields(cmdLine)
	for i, part := range parts {
		if part == "-L" && i+1 < len(parts) {
			// The next field should be "localport:host:remoteport"
			portMapping := parts[i+1]
			colonIdx := strings.Index(portMapping, ":")
			if colonIdx > 0 {
				portStr := portMapping[:colonIdx]
				port, err := strconv.Atoi(portStr)
				if err == nil {
					return port
				}
			}
		}
	}
	return 0
}

// isProcessRunning checks if a process with the given PID is running
func isProcessRunning(pid int) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid=")
	return cmd.Run() == nil
}

// RemoveTunnelState removes the state file for a tunnel on a given port
func RemoveTunnelState(localPort int) error {
	tunnelsDir, err := getTunnelsDir()
	if err != nil {
		return err
	}

	stateFile := filepath.Join(tunnelsDir, fmt.Sprintf("tunnel-%d.json", localPort))
	return os.Remove(stateFile)
}

// sanitizeDisplayName ensures the given string is a valid and safe display name by removing invalid characters and truncating the length.
func sanitizeDisplayName(s string) string {
	allowed := regexp.MustCompile(`[^A-Za-z0-9._+@-]`)
	clean := allowed.ReplaceAllString(s, "-")
	if len(clean) > 255 {
		clean = clean[:255]
	}
	if strings.Trim(clean, "-") == "" {
		clean = fmt.Sprintf("ocloud-%d", time.Now().Unix())
	}
	return clean
}

// waitForSessionActive polls the bastion session until it reaches ACTIVE or the context is cancelled.
// It mirrors the previous inline loops and keeps the small sleep after ACTIVE to ensure readiness.
func (s *Service) waitForSessionActive(ctx context.Context, sessionID string) error {
	for {
		getResp, err := s.bastionClient.GetSession(ctx, bastion.GetSessionRequest{SessionId: &sessionID})
		if err != nil {
			return fmt.Errorf("waiting for session ACTIVE: %w", err)
		}
		if getResp.Session.LifecycleState == bastion.SessionLifecycleStateActive {
			time.Sleep(waitPollInterval)
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitPollInterval):
		}
	}
}

// readPublicKey reads and returns the public key content from the given path.
func readPublicKey(publicKeyPath string) (string, error) {
	data, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("reading public key: %w", err)
	}
	return string(data), nil
}

// listActiveSessions returns ACTIVE session summaries for a bastion, sorted by time created desc.
func (s *Service) listActiveSessions(ctx context.Context, bastionID string) ([]bastion.SessionSummary, error) {
	lsReq := bastion.ListSessionsRequest{
		BastionId:             common.String(bastionID),
		SessionLifecycleState: bastion.ListSessionsSessionLifecycleStateActive,
		SortBy:                bastion.ListSessionsSortByTimecreated,
		SortOrder:             bastion.ListSessionsSortOrderDesc,
	}
	lsResp, err := s.bastionClient.ListSessions(ctx, lsReq)
	if err != nil {
		return nil, fmt.Errorf("listing bastion sessions: %w", err)
	}
	return lsResp.Items, nil
}

// EnsurePortForwardSession finds an ACTIVE bastion session targeting the given IP:port and matching the provided public key.
// If not found, it creates a new session and waits until it becomes ACTIVE, returning the session ID.
func (s *Service) EnsurePortForwardSession(ctx context.Context, bastionID, targetIP string, port int, publicKeyPath string) (string, error) {
	pubKey, err := readPublicKey(publicKeyPath)
	if err != nil {
		return "", err
	}

	// 1) Try to reuse an ACTIVE matching session
	items, err := s.listActiveSessions(ctx, bastionID)
	if err != nil {
		return "", err
	}
	for _, item := range items {
		if trd, ok := item.TargetResourceDetails.(bastion.PortForwardingSessionTargetResourceDetails); ok {
			if trd.TargetResourcePrivateIpAddress != nil && trd.TargetResourcePort != nil &&
				*trd.TargetResourcePrivateIpAddress == targetIP && *trd.TargetResourcePort == port {
				getResp, err := s.bastionClient.GetSession(ctx, bastion.GetSessionRequest{SessionId: item.Id})
				if err != nil {
					return "", fmt.Errorf("getting bastion session: %w", err)
				}
				if getResp.KeyDetails != nil && getResp.KeyDetails.PublicKeyContent != nil && *getResp.KeyDetails.PublicKeyContent == pubKey {
					return *item.Id, nil // Reuse
				}
			}
		}
	}

	// 2) Create a new session
	baseName := fmt.Sprintf("ocloud-%s-%d-%d", strings.ReplaceAll(targetIP, ".", "-"), port, time.Now().Unix())
	displayName := sanitizeDisplayName(baseName)
	createReq := bastion.CreateSessionRequest{
		CreateSessionDetails: bastion.CreateSessionDetails{
			BastionId: common.String(bastionID),
			TargetResourceDetails: bastion.CreatePortForwardingSessionTargetResourceDetails{
				TargetResourcePrivateIpAddress: common.String(targetIP),
				TargetResourcePort:             common.Int(port),
			},
			KeyDetails:          &bastion.PublicKeyDetails{PublicKeyContent: &pubKey},
			DisplayName:         common.String(displayName),
			SessionTtlInSeconds: common.Int(defaultTTL),
		},
	}
	crResp, err := s.bastionClient.CreateSession(ctx, createReq)
	if err != nil {
		return "", fmt.Errorf("creating bastion session: %w", err)
	}
	sessionID := *crResp.Id

	// 3) Wait for ACTIVE
	if err := s.waitForSessionActive(ctx, sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

// EnsureManagedSSHSession finds or creates a Managed SSH bastion session for the given target instance and returns the session ID.
func (s *Service) EnsureManagedSSHSession(ctx context.Context, bastionID, targetInstanceID, targetIP, osUser string, port int, publicKeyPath string, ttlSeconds int) (string, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultTTL
	}
	pubKey, err := readPublicKey(publicKeyPath)
	if err != nil {
		return "", err
	}

	//-------------------------Try to reuse an ACTIVE matching Managed SSH session--------------------------------------
	items, err := s.listActiveSessions(ctx, bastionID)
	if err != nil {
		return "", err
	}
	for _, item := range items {
		if trd, ok := item.TargetResourceDetails.(bastion.ManagedSshSessionTargetResourceDetails); ok {
			if trd.TargetResourceId != nil && trd.TargetResourcePrivateIpAddress != nil && trd.TargetResourcePort != nil && trd.TargetResourceOperatingSystemUserName != nil &&
				*trd.TargetResourceId == targetInstanceID && *trd.TargetResourcePrivateIpAddress == targetIP && *trd.TargetResourcePort == port && *trd.TargetResourceOperatingSystemUserName == osUser {
				getResp, err := s.bastionClient.GetSession(ctx, bastion.GetSessionRequest{SessionId: item.Id})
				if err != nil {
					return "", fmt.Errorf("getting bastion session: %w", err)
				}
				if getResp.KeyDetails != nil && getResp.KeyDetails.PublicKeyContent != nil && *getResp.KeyDetails.PublicKeyContent == pubKey {
					return *item.Id, nil
				}
			}
		}
	}

	//-----------------------------------------Create a new Managed SSH session-----------------------------------------
	baseName := fmt.Sprintf("ocloud-%s-%d-%d", strings.ReplaceAll(targetIP, ".", "-"), port, time.Now().Unix())
	displayName := sanitizeDisplayName(baseName)
	createReq := bastion.CreateSessionRequest{
		CreateSessionDetails: bastion.CreateSessionDetails{
			BastionId: common.String(bastionID),
			TargetResourceDetails: bastion.CreateManagedSshSessionTargetResourceDetails{
				TargetResourceId:                      common.String(targetInstanceID),
				TargetResourceOperatingSystemUserName: common.String(osUser),
				TargetResourcePort:                    common.Int(port),
				TargetResourcePrivateIpAddress:        common.String(targetIP),
			},
			KeyDetails:          &bastion.PublicKeyDetails{PublicKeyContent: &pubKey},
			DisplayName:         common.String(displayName),
			SessionTtlInSeconds: common.Int(ttlSeconds),
		},
	}
	crResp, err := s.bastionClient.CreateSession(ctx, createReq)
	if err != nil {
		return "", fmt.Errorf("creating bastion session: %w", err)
	}
	sessionID := *crResp.Id

	//------------------------------------------------Wait for ACTIVE---------------------------------------------------
	if err := s.waitForSessionActive(ctx, sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

// BuildManagedSSHCommand constructs the SSH command that uses ProxyCommand with the bastion Managed SSH session.
// It opens only a direct-tcpip channel on the bastion (accepted), while authenticating to bastion with the session OCID.
// The outer SSH connects to the target instance as targetUser@targetIP.
func BuildManagedSSHCommand(privateKeyPath, sessionID, region, targetIP, targetUser string) string {
	realm := "oraclecloud"
	parts := strings.Split(sessionID, ".")
	if len(parts) > 2 && strings.Contains(parts[2], "2") {
		realm = "oraclegovcloud"
	}
	proxy := fmt.Sprintf("ssh -i %s -W %%h:%%p -p 22 %s@host.bastion.%s.oci.%s.com", privateKeyPath, sessionID, region, realm)
	return fmt.Sprintf("ssh -i %s -o ProxyCommand=\"%s\" -p 22 %s@%s", privateKeyPath, proxy, targetUser, targetIP)
}

// BuildPortForwardArgs constructs SSH command arguments for establishing a secure port-forwarding tunnel.
// It handles path expansion for the private key, determines the correct realm domain, and formats connection options.
func BuildPortForwardArgs(privateKeyPath, sessionID, region, targetIP string, localPort, remotePort int) ([]string, error) {
	key, err := expandTilde(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("expand key path: %w", err)
	}

	// Decide realm domain based on OCID (oc2/oc3 => gov)
	realmDomain := "oraclecloud.com"
	if strings.Contains(sessionID, ".oc2.") || strings.Contains(sessionID, ".oc3.") {
		realmDomain = "oraclegovcloud.com"
	}

	bastionUser := fmt.Sprintf("%s@host.bastion.%s.oci.%s", sessionID, region, realmDomain)

	args := []string{
		"-i", key,
		"-o", "StrictHostKeyChecking=accept-new",
		// keepalives help the tunnel auto-detect dead links
		"-o", "ServerAliveInterval=30",
		"-o", "ServerAliveCountMax=3",
		"-N",
		"-L", fmt.Sprintf("%d:%s:%d", localPort, targetIP, remotePort),
		"-p", "22",
		bastionUser,
	}
	return args, nil
}

// expandTilde resolves paths beginning with "~" to the current user's home directory, returning the expanded path or an error.
func expandTilde(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, strings.TrimPrefix(p, "~")), nil
	}
	return p, nil
}

// SpawnDetached starts ssh in the background, detaches from your process, and returns its PID and log file path.
// localPort is used to generate a unique log file name.
func SpawnDetached(args []string, localPort int, targetIP string) (int, string, error) {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return 0, "", fmt.Errorf("ssh not found in PATH: %w", err)
	}

	// Create a log directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, "", fmt.Errorf("get home dir: %w", err)
	}
	profile := os.Getenv(flags.EnvKeyProfile)
	logDir := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCISessionsDirName, profile, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return 0, "", fmt.Errorf("create log dir: %w", err)
	}

	// Create a unique log file with a timestamp
	timestamp := time.Now().Format("20060102-150405")
	logfile := filepath.Join(logDir, fmt.Sprintf("ssh-tunnel-%d-%s.log", localPort, timestamp))

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 0, "", fmt.Errorf("open log file: %w", err)
	}
	defer f.Close()

	// Write header to a log file
	header := fmt.Sprintf("=== SSH Tunnel Started ===\nTimestamp: %s\nLocal Port: %d\nTarget: %s\nSSH Command: %s %s\n\n",
		time.Now().Format(time.RFC3339),
		localPort,
		targetIP,
		sshPath,
		strings.Join(args, " "))
	_, _ = f.WriteString(header)

	verboseArgs := append([]string{"-v"}, args...)

	cmd := exec.Command(sshPath, verboseArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // detach from our session/TTY
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		return 0, "", fmt.Errorf("start ssh: %w", err)
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return pid, logfile, nil
}

// WaitForListen wait until the localPort is listening (nice UX).
// Helps to avoid "connection refused" errors.
func WaitForListen(localPort int, timeout time.Duration) error {
	addr := fmt.Sprintf("127.0.0.1:%d", localPort)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		c, err := net.DialTimeout("tcp", addr, 400*time.Millisecond)
		if err == nil {
			_ = c.Close()
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("tunnel not up on %s after %s", addr, timeout)
}

// SpawnDetachedWithSudo starts ssh with sudo for privileged ports (below 1024).
// It runs in the background, detaches from your process, and returns its PID and log file path.
// The privateKeyPath is needed to explicitly specify the key path since sudo changes the user context.
func SpawnDetachedWithSudo(args []string, localPort int, targetIP string, privateKeyPath string) (int, string, error) {
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return 0, "", fmt.Errorf("sudo not found in PATH: %w", err)
	}

	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return 0, "", fmt.Errorf("ssh not found in PATH: %w", err)
	}

	// Create a log directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, "", fmt.Errorf("get home dir: %w", err)
	}
	profile := os.Getenv(flags.EnvKeyProfile)
	logDir := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCISessionsDirName, profile, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return 0, "", fmt.Errorf("create log dir: %w", err)
	}

	// Create a unique log file with a timestamp
	timestamp := time.Now().Format("20060102-150405")
	logfile := filepath.Join(logDir, fmt.Sprintf("ssh-tunnel-%d-%s.log", localPort, timestamp))

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 0, "", fmt.Errorf("open log file: %w", err)
	}
	defer f.Close()

	// Build the full command: sudo ssh -i <key> <args...>
	// We prepend the verbose flag and ensure the private key is explicitly specified
	verboseArgs := append([]string{"-v"}, args...)
	sudoArgs := append([]string{sshPath}, verboseArgs...)

	// Write header to a log file
	header := fmt.Sprintf("=== SSH Tunnel Started (with sudo) ===\nTimestamp: %s\nLocal Port: %d\nTarget: %s\nSSH Command: sudo %s %s\n\n",
		time.Now().Format(time.RFC3339),
		localPort,
		targetIP,
		sshPath,
		strings.Join(verboseArgs, " "))
	_, _ = f.WriteString(header)

	cmd := exec.Command(sudoPath, sudoArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // detach from our session/TTY
	cmd.Stdout = f
	cmd.Stderr = f
	// For sudo, we need to connect stdin for the password prompt
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return 0, "", fmt.Errorf("start sudo ssh: %w", err)
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return pid, logfile, nil
}
