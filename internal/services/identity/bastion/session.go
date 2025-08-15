package bastion

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	conf "github.com/rozdolsky33/ocloud/internal/config"
)

// Defaults used for session wait and ttl
var (
	waitPollInterval = 3 * time.Second
	defaultTTL       = 10800 // seconds (3 hours)
)

// DefaultSSHKeyPaths returns the default public and private key paths based on the active OCI CLI profile.
// It mirrors the sshConfig() defaults so callers can use keys without hardcoding paths.
func DefaultSSHKeyPaths() (publicKeyPath, privateKeyPath string) {
	homeDir, _ := conf.GetUserHomeDir()
	sessionDir := filepath.Join(homeDir, ".ssh")
	return filepath.Join(sessionDir, "id_rsa.pub"), filepath.Join(sessionDir, "id_rsa")
}

// sanitizeDisplayName ensures only allowed characters [A-Za-z0-9._+@-] and max length 255.
// Any disallowed rune is replaced with '-'. If a result is empty, returns a safe fallback.
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

// EnsurePortForwardSession finds an ACTIVE bastion session targeting the given IP:port and matching the provided public key.
// If not found, it creates a new session and waits until it becomes ACTIVE, returning the session ID.
func (s *Service) EnsurePortForwardSession(ctx context.Context, bastionID, targetIP string, port int, publicKeyPath string) (string, error) {
	pubKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("reading public key: %w", err)
	}
	pubKey := string(pubKeyData)

	// 1) Try to reuse an ACTIVE matching session
	lsReq := bastion.ListSessionsRequest{
		BastionId:             common.String(bastionID),
		SessionLifecycleState: bastion.ListSessionsSessionLifecycleStateActive,
		SortBy:                bastion.ListSessionsSortByTimecreated,
		SortOrder:             bastion.ListSessionsSortOrderDesc,
	}
	lsResp, err := s.bastionClient.ListSessions(ctx, lsReq)
	if err != nil {
		return "", fmt.Errorf("listing bastion sessions: %w", err)
	}
	for _, item := range lsResp.Items {
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
	pubKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("reading public key: %w", err)
	}
	pubKey := string(pubKeyData)

	// Try to reuse an ACTIVE matching Managed SSH session
	lsReq := bastion.ListSessionsRequest{
		BastionId:             common.String(bastionID),
		SessionLifecycleState: bastion.ListSessionsSessionLifecycleStateActive,
		SortBy:                bastion.ListSessionsSortByTimecreated,
		SortOrder:             bastion.ListSessionsSortOrderDesc,
	}
	lsResp, err := s.bastionClient.ListSessions(ctx, lsReq)
	if err != nil {
		return "", fmt.Errorf("listing bastion sessions: %w", err)
	}
	for _, item := range lsResp.Items {
		if trd, ok := item.TargetResourceDetails.(bastion.ManagedSshSessionTargetResourceDetails); ok {
			if trd.TargetResourceId != nil && trd.TargetResourcePrivateIpAddress != nil && trd.TargetResourcePort != nil && trd.TargetResourceOperatingSystemUserName != nil &&
				*trd.TargetResourceId == targetInstanceID && *trd.TargetResourcePrivateIpAddress == targetIP && *trd.TargetResourcePort == port && *trd.TargetResourceOperatingSystemUserName == osUser {
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

	// Create a new Managed SSH session
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

	// Wait for ACTIVE
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

// BuildPortForwardNohupCommand builds a nohup SSH command to run a local port forward via the bastion session in the background.
// It matches the example flags and routes localPort to targetIP:remotePort.
func BuildPortForwardNohupCommand(privateKeyPath, sessionID, region, targetIP string, localPort, remotePort int, logPath string) string {
	realm := "oraclecloud"
	parts := strings.Split(sessionID, ".")
	if len(parts) > 2 && strings.Contains(parts[2], "2") {
		realm = "oraclegovcloud"
	}
	bastionHost := fmt.Sprintf("%s@host.bastion.%s.oci.%s.com", sessionID, region, realm)
	if logPath == "" {
		logPath = "ssh-tunnel.log"
	}
	return fmt.Sprintf(
		"nohup ssh -i %s -o StrictHostKeyChecking=accept-new -o HostkeyAlgorithms=+ssh-rsa -o PubkeyAcceptedKeyTypes=+ssh-rsa -N -L %d:%s:%d -p 22 %s > %s 2>&1 &",
		privateKeyPath, localPort, targetIP, remotePort, bastionHost, logPath,
	)
}
