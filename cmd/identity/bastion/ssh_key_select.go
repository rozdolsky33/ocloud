package bastion

import (
	"bytes"
	"context"
	"crypto"
	xecdsa "crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"golang.org/x/crypto/ssh"
)

// SelectSSHKeyPair opens two TUIs to select a matching SSH public and private key.
// It returns ErrAborted if the user cancels either selection.
func SelectSSHKeyPair(ctx context.Context) (pubKey, privKey string, err error) {
	home := os.Getenv("HOME")
	if home == "" {
		if h, err := os.UserHomeDir(); err == nil {
			home = h
		}
	}
	startDir := filepath.Join(home, ".ssh")
	pk := NewSSHKeysModelBrowser("Choose Public Key", startDir, true)
	pProg := tea.NewProgram(pk, tea.WithContext(ctx))
	pRes, err := pProg.Run()
	if err != nil {
		return "", "", fmt.Errorf("public key selection TUI: %w", err)
	}
	pPick, ok := pRes.(SHHFilesModel)
	if !ok || pPick.Choice() == "" {
		return "", "", ErrAborted
	}
	pubKey = pPick.Choice()
	logger.CmdLogger.Info("selected public ssh key", "name", filepath.Base(pubKey), "path", pubKey)
	sk := NewSSHKeysModelBrowser("Choose Private Key", startDir, false)
	sProg := tea.NewProgram(sk, tea.WithContext(ctx))
	sRes, err := sProg.Run()
	if err != nil {
		return "", "", fmt.Errorf("private key selection TUI: %w", err)
	}
	sPick, ok := sRes.(SHHFilesModel)
	if !ok || sPick.Choice() == "" {
		return "", "", ErrAborted
	}
	privKey = sPick.Choice()
	logger.CmdLogger.Info("selected private ssh key", "name", filepath.Base(privKey), "path", privKey)

	expected := strings.TrimSuffix(pubKey, ".pub")
	if filepath.Base(privKey) != filepath.Base(expected) {
		return "", "", fmt.Errorf("selected private key %s does not match public key %s (expected private: %s)", privKey, pubKey, expected)
	}

	if err := validateSSHKeyPair(pubKey, privKey); err != nil {
		return "", "", err
	}

	return pubKey, privKey, nil
}

// validateSSHKeyPair ensures that:
// 1) the public key type is ssh-rsa, ssh-ed25519, or ecdsa-sha2-nistp{256,384,521},
// 2) the private key is RSA, ED25519, or ECDSA,
// 3) the derived public key from the private key matches the selected public key.
func validateSSHKeyPair(pubPath, privPath string) error {
	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return fmt.Errorf("read public key: %w", err)
	}
	pubKey, pubComment, _, _, pubErr := ssh.ParseAuthorizedKey(pubBytes)
	if pubErr != nil {
		return fmt.Errorf("parse public key: %w", pubErr)
	}
	pubType := pubKey.Type()
	if pubType != ssh.KeyAlgoRSA && pubType != ssh.KeyAlgoED25519 && pubType != ssh.KeyAlgoECDSA256 && pubType != ssh.KeyAlgoECDSA384 && pubType != ssh.KeyAlgoECDSA521 {
		return fmt.Errorf("unsupported public key type %s (allowed: ssh-rsa, ssh-ed25519, ecdsa-sha2-nistp256/384/521)%s", pubType, formatComment(pubComment))
	}

	privBytes, err := os.ReadFile(privPath)
	if err != nil {
		return fmt.Errorf("read private key: %w", err)
	}
	rawPriv, err := ssh.ParseRawPrivateKey(privBytes)
	if err != nil {
		if strings.Contains(err.Error(), "encrypted") {
			return errors.New("the selected private key is encrypted. please use an unencrypted key for this workflow or decrypt it temporarily")
		}
		return fmt.Errorf("parse private key: %w", err)
	}

	var allowed bool
	var derived ssh.PublicKey
	switch k := rawPriv.(type) {
	case *rsa.PrivateKey:
		allowed = true
		derived, err = ssh.NewPublicKey(&k.PublicKey)
	case ed25519.PrivateKey:
		allowed = true
		derived, err = ssh.NewPublicKey(k.Public())
	case *xecdsa.PrivateKey:
		allowed = true
		derived, err = ssh.NewPublicKey(&k.PublicKey)
	default:
		allowed = false
	}
	if !allowed {
		if s, ok := rawPriv.(crypto.Signer); ok {
			switch pk := s.Public().(type) {
			case ed25519.PublicKey:
				derived, err = ssh.NewPublicKey(pk)
				allowed = true
			case *xecdsa.PublicKey:
				derived, err = ssh.NewPublicKey(pk)
				allowed = true
			case *rsa.PublicKey:
				derived, err = ssh.NewPublicKey(pk)
				allowed = true
			default:
				return fmt.Errorf("unsupported private key type (allowed: RSA, ED25519, ECDSA)")
			}
		} else {
			return fmt.Errorf("unsupported private key type (allowed: RSA, ED25519, ECDSA)")
		}
	}
	if err != nil {
		return fmt.Errorf("derive public key from private: %w", err)
	}

	if isKeyAlgorithmMismatch(pubType, derived.Type()) {
		return fmt.Errorf("mismatch between public (%s) and private (%s) key algorithms", pubType, derived.Type())
	}

	if !bytes.Equal(pubKey.Marshal(), derived.Marshal()) {
		return fmt.Errorf("selected private key does not match the selected public key")
	}
	return nil
}

// isKeyAlgorithmMismatch returns true if the public and private key types do not match.
func isKeyAlgorithmMismatch(pubType, derivedType string) bool {
	return pubType != derivedType
}

// formatComment returns a comment string with a colon prefix if the comment is not empty.
func formatComment(c string) string {
	if c == "" {
		return ""
	}
	return ": " + c
}
