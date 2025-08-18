package bastion

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config"
	instancessvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
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

	clusters, _, _, err := okeService.List(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("list OKE clusters: %w", err)
	}
	if len(clusters) == 0 {
		fmt.Println("No OKE clusters found.")
		return nil
	}

	cm := NewOKEListModelFancy(clusters)
	cp := tea.NewProgram(cm, tea.WithContext(ctx))
	userSelection, err := cp.Run()
	if err != nil {
		return fmt.Errorf("OKE selection TUI: %w", err)
	}
	chosen, ok := userSelection.(ResourceListModel)
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

	// Session-type-specific handling for OKE
	switch sType {
	case TypeManagedSSH:
		// List compute instances and filter those starting with "oke"
		instService, err := instancessvc.NewService(appCtx)
		if err != nil {
			return fmt.Errorf("create instance service: %w", err)
		}
		instances, _, _, err := instService.List(ctx, 300, 0, true)
		if err != nil {
			return fmt.Errorf("list instances: %w", err)
		}
		filtered := make([]instancessvc.Instance, 0, len(instances))
		for _, it := range instances {
			if strings.HasPrefix(strings.ToLower(it.Name), "oke") {
				filtered = append(filtered, it)
			}
		}
		if len(filtered) == 0 {
			fmt.Println("No instances with name starting with 'oke' found.")
			return nil
		}

		// TUI selection among filtered instances
		im := NewInstanceListModelFancy(filtered)
		ip := tea.NewProgram(im, tea.WithContext(ctx))
		ires, err := ip.Run()
		if err != nil {
			return fmt.Errorf("instance selection TUI: %w", err)
		}
		chosenInstRes, ok := ires.(ResourceListModel)
		if !ok || chosenInstRes.Choice() == "" {
			return ErrAborted
		}
		var inst instancessvc.Instance
		for _, it := range filtered {
			if it.ID == chosenInstRes.Choice() {
				inst = it
				break
			}
		}

		pubKey, privKey, err := SelectSSHKeyPair(ctx)
		if err != nil {
			return err
		}

		// Reachability, for instance
		if ok, reason := svc.CanReach(ctx, b, inst.VcnID, inst.SubnetID); !ok {
			fmt.Println("Bastion cannot reach selected instance:", reason)
			return nil
		}

		region, regErr := appCtx.Provider.Region()
		if regErr != nil {
			return fmt.Errorf("get region: %w", regErr)
		}

		sshUser, err := util.PromptString("Enter SSH username", "opc")
		if err != nil {
			return fmt.Errorf("read ssh username: %w", err)
		}
		sessID, err := svc.EnsureManagedSSHSession(ctx, b.ID, inst.ID, inst.IP, sshUser, 22, pubKey, 0)
		if err != nil {
			return fmt.Errorf("ensure managed SSH: %w", err)
		}
		sshCmd := bastionSvc.BuildManagedSSHCommand(privKey, sessID, region, inst.IP, sshUser)
		fmt.Printf("\nExecuting: %s\n\n", sshCmd)
		return bastionSvc.RunShell(ctx, appCtx.Stdout, appCtx.Stderr, sshCmd)

	case TypePortForwarding:
		// For PF: resolve endpoint -> private IP, then tunnel to 6443
		candidates := []string{}
		if h := util.ExtractHostname(cluster.PrivateEndpoint); h != "" {
			candidates = append(candidates, h)
		}
		if h := util.ExtractHostname(cluster.KubernetesEndpoint); h != "" {
			candidates = append(candidates, h)
		}
		if len(candidates) == 0 {
			return fmt.Errorf("could not determine OKE API host from endpoints: kube=%q private=%q",
				cluster.KubernetesEndpoint, cluster.PrivateEndpoint)
		}

		var targetIP string
		var lastErr error
		for _, host := range candidates {
			ip, err := util.ResolveHostToIP(ctx, host) // ctx-aware DNS
			if err == nil {
				targetIP = ip
				break
			}
			lastErr = err
		}
		if targetIP == "" {
			return fmt.Errorf("resolve OKE API endpoint to private IP: %v", lastErr)
		}

		pubKey, privKey, err := SelectSSHKeyPair(ctx)
		if err != nil {
			return err
		}

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

		if util.IsLocalTCPPortInUse(port) {
			return fmt.Errorf("local port %d is already in use on 127.0.0.1; choose another port", port)
		}

		// Ensure kubeconfig only if missing; prompt the user before creating/merging
		exists, err := okesvc.KubeconfigExistsForOKE(cluster, region, config.GetOCIProfile())
		if err != nil {
			return fmt.Errorf("check kubeconfig: %w", err)
		}
		if !exists {
			question := "Kubeconfig for this OKE cluster was not found in ~/.kube/config. Create and merge it now?"
			if util.PromptYesNo(question) {
				if err := okesvc.EnsureKubeconfigForOKE(cluster, region, config.GetOCIProfile(), port); err != nil {
					return fmt.Errorf("ensure kubeconfig: %w", err)
				}
			} else {
				fmt.Println("Skipping kubeconfig creation for this OKE cluster.")
			}
		}

		localPort := port
		logFile := fmt.Sprintf("~/.oci/.ocloud/ssh-tunnel-%d.log", localPort)
		sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, targetIP, localPort, okeTargetPort)

		if err != nil {
			return fmt.Errorf("build args: %w", err)
		}

		pid, err := bastionSvc.SpawnDetached(sshTunnelArgs, "/tmp/ssh-tunnel.log")
		if err != nil {
			return fmt.Errorf("spawn detached: %w", err)
		}
		log.Printf("spawned tunnel pid=%d", pid)

		// (optional)
		if err := bastionSvc.WaitForListen(okeTargetPort, 5*time.Second); err != nil {
			log.Printf("warning: %v", err)
		}

		fmt.Printf("\nStarting background OKE API tunnel: %s\n\n", sshTunnelArgs)

		fmt.Printf("SSH tunnel to OKE API started. Access: https://127.0.0.1:%d (kube-apiserver)\nLogs: %s\n",
			localPort, logFile)
		return nil
	default:
		return fmt.Errorf("unsupported session type: %s", sType)
	}
}
