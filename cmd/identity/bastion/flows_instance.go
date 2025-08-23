package bastion

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociInst "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
	instSvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// connectInstance runs the flow for an Instance target.
func connectInstance(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

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

	instances, _, _, err := instService.List(ctx, 300, 0)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}

	if len(instances) == 0 {
		logger.Logger.Info("No instances found.")
		return nil
	}

	// TUI selection
	im := NewInstanceListModelFancy(instances)
	ip := tea.NewProgram(im, tea.WithContext(ctx))
	ires, err := ip.Run()
	if err != nil {
		return fmt.Errorf("instance selection TUI: %w", err)
	}
	chosen, ok := ires.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return ErrAborted
	}

	pubKey, privKey, err := SelectSSHKeyPair(ctx)
	if err != nil {
		return err
	}

	var inst instSvc.Instance
	for _, it := range instances {
		if it.OCID == chosen.Choice() {
			inst = it
			break
		}
	}

	if ok, reason := svc.CanReach(ctx, b, inst.VcnID, inst.SubnetID); !ok {
		logger.Logger.Info("Bastion cannot reach selected instance", "reason", reason)
		return nil
	}

	logger.Logger.Info("Validated session on Bastion to Instance", "session_type", sType, "bastion_name", b.Name, "bastion_id", b.ID, "instance_name", inst.DisplayName)

	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	switch sType {
	case TypeManagedSSH:
		sshUser, err := util.PromptString("Enter SSH username", "opc")
		if err != nil {
			return fmt.Errorf("read ssh username: %w", err)
		}
		sessID, err := svc.EnsureManagedSSHSession(ctx, b.ID, inst.OCID, inst.PrimaryIP, sshUser, 22, pubKey, 0)
		if err != nil {
			return fmt.Errorf("ensure managed SSH: %w", err)
		}
		sshCmd := bastionSvc.BuildManagedSSHCommand(privKey, sessID, region, inst.PrimaryIP, sshUser)
		logger.Logger.Info("Executing", "command", sshCmd)
		return bastionSvc.RunShell(ctx, appCtx.Stdout, appCtx.Stderr, sshCmd)
	case TypePortForwarding:
		defaultPort := 5901
		port, err := util.PromptPort("Enter port to forward (local:target)", defaultPort)
		if err != nil {
			return fmt.Errorf("read port: %w", err)
		}
		sessID, err := svc.EnsurePortForwardSession(ctx, b.ID, inst.PrimaryIP, port, pubKey)
		if err != nil {
			return fmt.Errorf("ensure port forward: %w", err)
		}
		logFile := fmt.Sprintf("~/.oci/.ocloud/ssh-tunnel-%d.log", port)
		sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, inst.PrimaryIP, port, port)
		if err != nil {
			return fmt.Errorf("build args: %w", err)
		}

		logger.Logger.Info("Starting background tunnel", "args", sshTunnelArgs)
		pid, err := bastionSvc.SpawnDetached(sshTunnelArgs, "/tmp/ssh-tunnel.log")

		if err != nil {
			return fmt.Errorf("spawn detached: %w", err)
		}
		logger.Logger.V(1).Info("spawned tunnel", "pid", pid)

		if err := bastionSvc.WaitForListen(defaultPort, 5*time.Second); err != nil {
			logger.Logger.Error(err, "warning")
		}

		logger.Logger.Info("Starting background tunnel", "args", sshTunnelArgs)

		logger.Logger.Info("SSH tunnel started in background", "logs", logFile)
		return nil
	default:
		return fmt.Errorf("unsupported session type: %s", sType)
	}
}
