package bastion

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// SelectSSHKeyPair opens two TUIs to select a matching SSH public and private key.
// It returns ErrAborted if the user cancels either selection.
func SelectSSHKeyPair(ctx context.Context) (pubKey, privKey string, err error) {
	pk := NewSSHKeysModelFancyList("Choose Public Key", util.DefaultPublicSSHKey())
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

	sk := NewSSHKeysModelFancyList("Choose Private Key", util.DefaultPrivateSSHKeys())
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
