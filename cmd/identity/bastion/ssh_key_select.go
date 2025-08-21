package bastion

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	return pubKey, privKey, nil
}
