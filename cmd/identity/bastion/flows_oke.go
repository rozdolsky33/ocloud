package bastion

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	okesvc "github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// connectOKE runs the OKE target flow. For Port-Forwarding, it resolves the API
// endpoint to a private IP and spawns a tunnel (nohup style).
func connectOKE(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

	okeService, err := okesvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("create OKE service: %w", err)
	}

	clusters, _, _, err := okeService.List(ctx, 50, 0)
	if err != nil {
		return fmt.Errorf("list OKE clusters: %w", err)
	}
	if len(clusters) == 0 {
		fmt.Println("No OKE clusters found.")
		return nil
	}

	cm := NewOKEListModelFancy(clusters)
	cp := tea.NewProgram(cm, tea.WithContext(ctx))
	cres, err := cp.Run()
	if err != nil {
		return fmt.Errorf("OKE selection TUI: %w", err)
	}
	chosen, ok := cres.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return ErrAborted
	}

	var cluster okesvc.Cluster
	for _, c := range clusters {
		if c.ID == chosen.Choice() {
			cluster = c
			break
		}
	}

	// Reachability check (VcnID is sufficient for cluster-level)
	if ok, reason := svc.CanReach(ctx, b, cluster.VcnID, ""); !ok {
		fmt.Println("Bastion cannot reach selected OKE cluster:", reason)
		return nil
	}

	fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to OKE cluster %s.\n",
		sType, b.Name, b.ID, cluster.Name)

	if sType != TypePortForwarding {
		// Managed SSH doesn't apply to OKE; just acknowledge prep.
		return nil
	}

	// For PF: resolve endpoint -> private IP, then tunnel to 6443
	candidates := []string{}
	if h := extractHostname(cluster.PrivateEndpoint); h != "" {
		candidates = append(candidates, h)
	}
	if h := extractHostname(cluster.KubernetesEndpoint); h != "" {
		candidates = append(candidates, h)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("could not determine OKE API host from endpoints: kube=%q private=%q",
			cluster.KubernetesEndpoint, cluster.PrivateEndpoint)
	}

	var targetIP string
	var lastErr error
	for _, host := range candidates {
		ip, err := resolveHostToIP(ctx, host) // ctx-aware DNS
		if err == nil {
			targetIP = ip
			break
		}
		lastErr = err
	}
	if targetIP == "" {
		return fmt.Errorf("resolve OKE API endpoint to private IP: %v", lastErr)
	}

	pubKey, privKey := bastionSvc.DefaultSSHKeyPaths()
	okeTargetPort := 6443
	sessID, err := svc.EnsurePortForwardSession(ctx, b.ID, targetIP, okeTargetPort, pubKey)
	if err != nil {
		return fmt.Errorf("ensure port forward: %w", err)
	}

	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	port, err := util.PromptPort("Enter port to forward (local:target)", okeTargetPort)
	if err != nil {
		return fmt.Errorf("read port: %w", err)
	}
	localPort := port
	logFile := fmt.Sprintf("~/.oci/.ocloud/ssh-tunnel-%d.log", localPort)
	sshCmd := bastionSvc.BuildPortForwardNohupCommand(privKey, sessID, region, targetIP, localPort, okeTargetPort, logFile)

	fmt.Printf("\nStarting background OKE API tunnel: %s\n\n", sshCmd)
	if err := RunShell(ctx, appCtx.Stdout, appCtx.Stderr, sshCmd); err != nil {
		return fmt.Errorf("start SSH tunnel: %w", err)
	}
	fmt.Printf("SSH tunnel to OKE API started. Access: https://127.0.0.1:%d (kube-apiserver)\nLogs: %s\n",
		localPort, logFile)
	return nil
}
