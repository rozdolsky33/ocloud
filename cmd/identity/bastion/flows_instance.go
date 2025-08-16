package bastion

import (
	"context"
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	instancessvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// connectInstance runs the flow for an Instance target. It validates reachability
// and creates either Managed SSH or Port-Forwarding session, spawning SSH accordingly.
func connectInstance(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

	instService, err := instancessvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("create instance service: %w", err)
	}

	instances, _, _, err := instService.List(ctx, 300, 0, true)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	if len(instances) == 0 {
		fmt.Println("No instances found.")
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

	// Find selection
	var inst instancessvc.Instance
	for _, it := range instances {
		if it.ID == chosen.Choice() {
			inst = it
			break
		}
	}

	// Reachability
	if ok, reason := svc.CanReach(ctx, b, inst.VcnID, inst.SubnetID); !ok {
		fmt.Println("Bastion cannot reach selected instance:", reason)
		return nil
	}

	fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to Instance %s.\n",
		sType, b.Name, b.ID, inst.Name)

	pubKey, privKey := bastionSvc.DefaultSSHKeyPaths()
	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	switch sType {
	case TypeManagedSSH:
		sshUser := "opc"
		sessID, err := svc.EnsureManagedSSHSession(ctx, b.ID, inst.ID, inst.IP, sshUser, 22, pubKey, 0)
		if err != nil {
			return fmt.Errorf("ensure managed SSH: %w", err)
		}
		sshCmd := bastionSvc.BuildManagedSSHCommand(privKey, sessID, region, inst.IP, sshUser)
		fmt.Printf("\nExecuting: %s\n\n", sshCmd)
		return bastionSvc.RunShell(ctx, appCtx.Stdout, appCtx.Stderr, sshCmd)
		//TODO: verify approach
	case TypePortForwarding:
		defaultPort := 5901
		port, err := util.PromptPort("Enter port to forward (local:target)", defaultPort)
		if err != nil {
			return fmt.Errorf("read port: %w", err)
		}
		sessID, err := svc.EnsurePortForwardSession(ctx, b.ID, inst.IP, port, pubKey)
		if err != nil {
			return fmt.Errorf("ensure port forward: %w", err)
		}
		logFile := fmt.Sprintf("~/.oci/.ocloud/ssh-tunnel-%d.log", port)
		sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, inst.IP, port, port)
		if err != nil {
			return fmt.Errorf("build args: %w", err)
		}

		fmt.Printf("\nStarting background tunnel: %s\n\n", sshTunnelArgs)
		pid, err := bastionSvc.SpawnDetached(sshTunnelArgs, "/tmp/ssh-tunnel.log")

		if err != nil {
			return fmt.Errorf("spawn detached: %w", err)
		}
		log.Printf("spawned tunnel pid=%d", pid)

		// (optional)
		if err := bastionSvc.WaitForListen(defaultPort, 5*time.Second); err != nil {
			log.Printf("warning: %v", err)
		}

		fmt.Printf("\nStarting background tunnel: %s\n\n", sshTunnelArgs)

		fmt.Printf("SSH tunnel started in background. Logs: %s\n", logFile)
		return nil
	default:
		return fmt.Errorf("unsupported session type: %s", sType)
	}
}
