package bastion

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

// SelectSSHKeyPair opens two TUIs to select a matching SSH public and private key.
// It returns ErrAborted if the user cancels either selection.
func SelectSSHKeyPair(ctx context.Context) (pubKey, privKey string, err error) {
	// Start at ~/.ssh for browsing
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

	// You can browse to choose a private key; we'll validate the match afterward.
	// Allow browsing non-.pub files for private key selection
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
// 1) the public key type is ssh-rsa or ssh-ed25519,
// 2) the private key is RSA or ED25519,
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
	if pubType != ssh.KeyAlgoRSA && pubType != ssh.KeyAlgoED25519 {
		return fmt.Errorf("unsupported public key type %s (allowed: ssh-rsa, ssh-ed25519)%s", pubType, formatComment(pubComment))
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
	default:
		allowed = false
	}
	if !allowed {
		return fmt.Errorf("unsupported private key type (allowed: RSA, ED25519)")
	}
	if err != nil {
		return fmt.Errorf("derive public key from private: %w", err)
	}

	if (pubType == ssh.KeyAlgoRSA && derived.Type() != ssh.KeyAlgoRSA) || (pubType == ssh.KeyAlgoED25519 && derived.Type() != ssh.KeyAlgoED25519) {
		return fmt.Errorf("mismatch between public (%s) and private (%s) key algorithms", pubType, derived.Type())
	}

	if !bytes.Equal(pubKey.Marshal(), derived.Marshal()) {
		return fmt.Errorf("selected private key does not match the selected public key")
	}
	return nil
}

func formatComment(c string) string {
	if c == "" {
		return ""
	}
	return ": " + c
}
