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

	// Filter private keys to match the selected public key basename (without .pub)
	privList := util.DefaultPrivateSSHKeys()
	expectedBase := filepath.Base(strings.TrimSuffix(pubKey, ".pub"))
	filteredPriv := make([]string, 0, len(privList))
	for _, p := range privList {
		if filepath.Base(p) == expectedBase {
			filteredPriv = append(filteredPriv, p)
		}
	}
	var privOptions []string
	if len(filteredPriv) > 0 {
		privOptions = filteredPriv
	} else {
		privOptions = privList
	}
	sk := NewSSHKeysModelFancyList("Choose Private Key", privOptions)
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
