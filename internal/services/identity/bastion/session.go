package bastion

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
)

// Defaults used for session wait and ttl
var (
	waitPollInterval = 3 * time.Second
	defaultTTL       = 10800 // seconds (3 hours)
)

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
	// Expand "~" if present
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

// SpawnDetached starts ssh in the background, detaches from your process, and returns its PID.
func SpawnDetached(args []string, logfile string) (int, error) {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return 0, fmt.Errorf("ssh not found in PATH: %w", err)
	}

	// Ensure log dir exists
	if err := os.MkdirAll(filepath.Dir(logfile), 0o755); err != nil {
		return 0, fmt.Errorf("create log dir: %w", err)
	}
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 0, fmt.Errorf("open log file: %w", err)
	}
	defer f.Close()

	cmd := exec.Command(sshPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // detach from our session/TTY
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("start ssh: %w", err)
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return pid, nil
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
