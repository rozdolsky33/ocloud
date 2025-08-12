package bastion

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	conf "github.com/rozdolsky33/ocloud/internal/config"
	cflags "github.com/rozdolsky33/ocloud/internal/config/flags"
)

// Defaults used for session wait and ttl
var (
	waitPollInterval = 3 * time.Second
	defaultTTL       = 3600 // seconds
)

// DefaultSSHKeyPaths returns the default public and private key paths based on the active OCI CLI profile.
// It mirrors the sshConfig() defaults so callers can use keys without hardcoding paths.
func DefaultSSHKeyPaths() (publicKeyPath, privateKeyPath string) {
	// Build ~/.oci/sessions/<PROFILE>
	homeDir, _ := conf.GetUserHomeDir()
	profile := conf.GetOCIProfile()
	sessionDir := filepath.Join(homeDir, cflags.OCIConfigDirName, cflags.OCISessionsDirName, profile)
	return filepath.Join(sessionDir, "oci_api_key_public.pem"), filepath.Join(sessionDir, "oci_api_key.pem")
}

// EnsurePortForwardSession finds an ACTIVE bastion session targeting the given IP:port and matching the provided public key.
// If not found, it creates a new session and waits until it becomes ACTIVE, returning the session ID.
func (s *Service) EnsurePortForwardSession(ctx context.Context, bastionID, targetIP string, port int, publicKeyPath string, ttlSeconds int) (string, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultTTL
	}
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
	createReq := bastion.CreateSessionRequest{
		CreateSessionDetails: bastion.CreateSessionDetails{
			BastionId: common.String(bastionID),
			TargetResourceDetails: bastion.CreatePortForwardingSessionTargetResourceDetails{
				TargetResourcePrivateIpAddress: common.String(targetIP),
				TargetResourcePort:             common.Int(port),
			},
			KeyDetails:          &bastion.PublicKeyDetails{PublicKeyContent: &pubKey},
			DisplayName:         common.String(fmt.Sprintf("ocloud-%s-%d", targetIP, time.Now().Unix())),
			SessionTtlInSeconds: common.Int(ttlSeconds),
		},
	}
	crResp, err := s.bastionClient.CreateSession(ctx, createReq)
	if err != nil {
		return "", fmt.Errorf("creating bastion session: %w", err)
	}
	sessionID := *crResp.Id

	// 3) Wait for ACTIVE
	for {
		getResp, err := s.bastionClient.GetSession(ctx, bastion.GetSessionRequest{SessionId: &sessionID})
		if err != nil {
			return "", fmt.Errorf("waiting for session ACTIVE: %w", err)
		}
		if getResp.Session.LifecycleState == bastion.SessionLifecycleStateActive {
			// Extra little delay as immediate connections can sometimes fail even after ACTIVE
			time.Sleep(waitPollInterval)
			break
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(waitPollInterval):
		}
	}
	return sessionID, nil
}

// BuildProxySSHCommand constructs the SSH command (as string) using ProxyCommand with the bastion session ID as username.
// The realm is inferred from the session OCID: if it contains "2" in the realm segment, use oraclegovcloud; else oraclecloud.
func BuildProxySSHCommand(privateKeyPath, sessionID, region, targetIP string) string {
	realm := "oraclecloud"
	// OCID format: ocid1.bastionsession.oc[1|2].... Split by '.' and inspect the 3rd token (index 2)
	parts := strings.Split(sessionID, ".")
	if len(parts) > 2 && strings.Contains(parts[2], "2") {
		realm = "oraclegovcloud"
	}
	proxy := fmt.Sprintf("ssh -i %s -W %%h:%%p -p 22 %s@host.bastion.%s.oci.%s.com", privateKeyPath, sessionID, region, realm)
	return fmt.Sprintf("ssh -i %s -o ProxyCommand=\"%s\" -p 22 opc@%s", privateKeyPath, proxy, targetIP)
}
