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
	ocilb "github.com/rozdolsky33/ocloud/internal/oci/network/loadbalancer"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	lbSvc "github.com/rozdolsky33/ocloud/internal/services/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// connectLoadBalancer runs the Load Balancer target flow for port forwarding.
func connectLoadBalancer(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

	// Only Port-Forwarding is supported for Load Balancers
	if sType != TypePortForwarding {
		logger.Logger.Info("Only Port-Forwarding sessions are supported for Load Balancer connections")
		return fmt.Errorf("only Port-Forwarding sessions are supported for Load Balancer connections")
	}

	// Create Load Balancer clients and service
	lbClient, err := oci.NewLoadBalancerClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating load balancer client: %w", err)
	}
	nwClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}
	certsClient, err := oci.NewCertificatesManagementClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating certificates management client: %w", err)
	}
	adapter := ocilb.NewAdapter(lbClient, nwClient, certsClient)
	lbService := lbSvc.NewService(adapter, appCtx)

	// Fetch load balancers
	allLBs, err := lbService.ListLoadBalancers(ctx)
	if err != nil {
		return fmt.Errorf("list load balancers: %w", err)
	}

	// Filter to only private load balancers (bastion can only reach private IPs)
	var lbs []lbSvc.LoadBalancer
	for _, lb := range allLBs {
		if strings.ToLower(lb.Type) == "private" {
			lbs = append(lbs, lb)
		}
	}

	if len(lbs) == 0 {
		logger.Logger.Info("No private Load Balancers found. Note: Only private Load Balancers can be accessed via bastion port forwarding.")
		return nil
	}

	// Display TUI for load balancer selection
	lm := NewLoadBalancerListModelFancy(lbs)
	lp := tea.NewProgram(lm, tea.WithContext(ctx))
	lres, err := lp.Run()
	if err != nil {
		return fmt.Errorf("load balancer selection TUI: %w", err)
	}
	chosen, ok := lres.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return ErrAborted
	}

	// Find the selected load balancer
	var lb lbSvc.LoadBalancer
	for _, l := range lbs {
		if l.OCID == chosen.Choice() {
			lb = l
			break
		}
	}

	// Check bastion reachability to load balancer's VCN
	subnetID := ""
	if len(lb.Subnets) > 0 {
		subnetID = lb.Subnets[0]
	}
	_, reason := svc.CanReach(ctx, b, lb.VcnID, subnetID)
	logger.Logger.Info("Reachability to Load Balancer cannot be automatically verified", "reason", reason)
	logger.Logger.Info("Validated session on Bastion to Load Balancer",
		"session_type", sType,
		"bastion_name", b.DisplayName,
		"bastion_id", b.OCID,
		"lb_name", lb.Name,
		"lb_ip_addresses", lb.IPAddresses,
		"lb_subnets", lb.Subnets,
		"lb_vcn_id", lb.VcnID)

	// Get SSH key pair
	pubKey, privKey, err := SelectSSHKeyPair(ctx)
	if err != nil {
		return err
	}

	// Get region
	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	// Determine target IP from load balancer
	// Note: IPAddresses may contain suffix like " (private)" or " (public)" from mapping
	if len(lb.IPAddresses) == 0 {
		return fmt.Errorf("no IP addresses found for load balancer %s", lb.Name)
	}
	targetIP := extractIPAddress(lb.IPAddresses[0])

	// The LB target port (what the LB is listening to on) - typically 443
	lbTargetPort := 443

	// The default local port is 8443 to avoid sudo requirement
	// User can choose 443 if they want to match the LB port (requires sudo)
	defaultLocalPort := 8443

	// Prompt for local port
	localPort, err := util.PromptPort("Enter local port to forward", defaultLocalPort)
	if err != nil {
		return fmt.Errorf("read port: %w", err)
	}

	// Check if the local port is already in use
	if util.IsLocalTCPPortInUse(localPort) {
		return fmt.Errorf("local port %d is already in use on 127.0.0.1; choose another port", localPort)
	}

	// Create a port forwarding session to the LB's target port
	logger.Logger.Info("Creating port forwarding session",
		"bastion_id", b.OCID,
		"target_ip", targetIP,
		"lb_target_port", lbTargetPort,
		"local_port", localPort)
	sessID, err := svc.EnsurePortForwardSession(ctx, b.OCID, targetIP, lbTargetPort, pubKey)
	if err != nil {
		return fmt.Errorf("ensure port forward: %w", err)
	}

	// Build SSH tunnel arguments: localPort -> targetIP:lbTargetPort
	sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, targetIP, localPort, lbTargetPort)
	if err != nil {
		return fmt.Errorf("build args: %w", err)
	}

	logger.Logger.Info("Starting SSH tunnel",
		"local_port", localPort,
		"target", fmt.Sprintf("%s:%d", targetIP, lbTargetPort),
		"lb_name", lb.Name)

	// Spawn detached in background
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

	logger.Logger.Info("SSH tunnel to Load Balancer running",
		"access", fmt.Sprintf("https://127.0.0.1:%d", localPort),
		"lb_name", lb.Name,
		"logs", logFile)
	return nil
}

// extractIPAddress extracts just the IP address from a string that may contain
// a suffix like " (private)" or " (public)" added by the mapping layer.
func extractIPAddress(ipWithSuffix string) string {
	// Handle formats like "217.142.42.8 (private)" or "155.248.24.108 (public)"
	if idx := strings.Index(ipWithSuffix, " "); idx > 0 {
		return ipWithSuffix[:idx]
	}
	return ipWithSuffix
}
