package bastion

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociInst "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
	ociOke "github.com/rozdolsky33/ocloud/internal/oci/compute/oke"
	instSvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	okeSvc "github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// connectOKE runs the OKE target flow.
func connectOKE(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

	containerEngineClient, err := oci.NewContainerEngineClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating container engine client: %w", err)
	}
	okeAdapter := ociOke.NewAdapter(containerEngineClient)
	okeService := okeSvc.NewService(okeAdapter, appCtx.Logger, appCtx.CompartmentID)

	clusters, _, _, err := okeService.FetchPaginatedClusters(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("list OKE clusters: %w", err)
	}
	if len(clusters) == 0 {
		logger.Logger.Info("No OKE clusters found.")
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

	var cluster okeSvc.Cluster
	for _, c := range clusters {
		if c.OCID == chosen.Choice() {
			cluster = c
			break
		}
	}

	if ok, reason := svc.CanReach(ctx, b, cluster.VcnOCID, ""); !ok {
		logger.Logger.Info("Bastion cannot reach selected OKE cluster", "reason", reason)
		return nil
	}

	logger.Logger.Info("Validated session on Bastion to OKE cluster", "session_type", sType, "bastion_name", b.DisplayName, "bastion_id", b.OCID, "cluster_name", cluster.DisplayName)

	switch sType {
	case TypeManagedSSH:
		computeClient, err := oci.NewComputeClient(appCtx.Provider)
		if err != nil {
			return fmt.Errorf("creating compute client: %w", err)
		}
		networkClient, err := oci.NewNetworkClient(appCtx.Provider)
		if err != nil {
			return fmt.Errorf("creating network client: %w", err)
		}
		instanceAdapter := ociInst.NewAdapter(computeClient, networkClient)
		instService := instSvc.NewService(instanceAdapter, appCtx.Logger, appCtx.CompartmentID)

		instances, _, _, err := instService.FetchPaginatedInstances(ctx, 300, 0)
		if err != nil {
			return fmt.Errorf("list instances: %w", err)
		}
		filtered := make([]instSvc.Instance, 0, len(instances))
		for _, it := range instances {
			if strings.HasPrefix(strings.ToLower(it.DisplayName), "oke") {
				filtered = append(filtered, it)
			}
		}
		if len(filtered) == 0 {
			logger.Logger.Info("No instances with name starting with 'oke' found.")
			return nil
		}

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
		var inst instSvc.Instance
		for _, it := range filtered {
			if it.OCID == chosenInstRes.Choice() {
				inst = it
				break
			}
		}

		pubKey, privKey, err := SelectSSHKeyPair(ctx)
		if err != nil {
			return err
		}

		if ok, reason := svc.CanReach(ctx, b, inst.VcnID, inst.SubnetID); !ok {
			logger.Logger.Info("Bastion cannot reach selected instance", "reason", reason)
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
		sessID, err := svc.EnsureManagedSSHSession(ctx, b.OCID, inst.OCID, inst.PrimaryIP, sshUser, 22, pubKey, 0)
		if err != nil {
			return fmt.Errorf("ensure managed SSH: %w", err)
		}
		sshCmd := bastionSvc.BuildManagedSSHCommand(privKey, sessID, region, inst.PrimaryIP, sshUser)
		logger.Logger.Info("Executing", "command", sshCmd)
		return bastionSvc.RunShell(ctx, appCtx.Stdout, appCtx.Stderr, sshCmd)

	case TypePortForwarding:
		candidates := []string{}
		if h := util.ExtractHostname(cluster.PrivateEndpoint); h != "" {
			candidates = append(candidates, h)
		}
		if h := util.ExtractHostname(cluster.PublicEndpoint); h != "" {
			candidates = append(candidates, h)
		}
		if len(candidates) == 0 {
			return fmt.Errorf("could not determine OKE API host from endpoints: kube=%q private=%q",
				cluster.PublicEndpoint, cluster.PrivateEndpoint)
		}

		var targetIP string
		var lastErr error
		for _, host := range candidates {
			ip, err := util.ResolveHostToIP(ctx, host)
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
		sessID, err := svc.EnsurePortForwardSession(ctx, b.OCID, targetIP, okeTargetPort, pubKey)
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

		exists, err := okeSvc.KubeconfigExistsForOKE(cluster, region)
		if err != nil {
			return fmt.Errorf("check kubeconfig: %w", err)
		}
		if !exists {
			question := "Kubeconfig for this OKE cluster was not found in ~/.kube/config. Create and merge it now?"
			if util.PromptYesNo(question) {
				if err := okeSvc.EnsureKubeconfigForOKE(cluster, region, port); err != nil {
					return fmt.Errorf("ensure kubeconfig: %w", err)
				}
			} else {
				logger.Logger.Info("Skipping kubeconfig creation for this OKE cluster.")
			}
		}

		localPort := port
		sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, targetIP, localPort, okeTargetPort)

		if err != nil {
			return fmt.Errorf("build args: %w", err)
		}

		pid, logFile, err := bastionSvc.SpawnDetached(sshTunnelArgs, localPort, targetIP)
		if err != nil {
			return fmt.Errorf("spawn detached: %w", err)
		}
		logger.Logger.V(logger.Debug).Info("spawned tunnel", "pid", pid)

		// Save tunnel state for tracking
		tunnelInfo := bastionSvc.TunnelInfo{
			PID:       pid,
			LocalPort: localPort,
			TargetIP:  targetIP,
			StartedAt: time.Now(),
			LogFile:   logFile,
		}
		if err := bastionSvc.SaveTunnelState(tunnelInfo); err != nil {
			logger.Logger.Error(err, "failed to save tunnel state")
		}

		logger.Logger.Info("SSH tunnel process started, waiting for connection to be ready...")
		if err := bastionSvc.WaitForListen(localPort, 30*time.Second); err != nil {
			logger.Logger.Info("Tunnel verification timed out, but the tunnel may still be establishing in the background", "port", localPort)
			logger.Logger.Info("Check the tunnel status and logs if you experience connection issues")
		} else {
			logger.Logger.Info("Tunnel is ready and accepting connections")
		}

		logger.Logger.Info("SSH tunnel to OKE API running", "access", fmt.Sprintf("https://127.0.0.1:%d (kube-apiserver)", localPort), "logs", logFile)
		return nil
	default:
		return fmt.Errorf("unsupported session type: %s", sType)
	}
}
