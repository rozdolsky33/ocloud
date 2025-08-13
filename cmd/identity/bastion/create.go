package bastion

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/rozdolsky33/ocloud/internal/app"
	instancessvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	okesvc "github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	autonomousdbsvc "github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/spf13/cobra"
)

func NewCreateCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create",
		Aliases:       []string{"c"},
		Short:         "Create a new Bastion",
		Long:          "Create a new Bastion by selecting from available options",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCreateCommand(cmd, appCtx)
		},
	}

	return cmd
}

// RunCreateCommand executes the create command with an interactive TUI
// It implements a two-step selection process:
// 1. First, the user selects a type (Bastion or Session)
// 2. Then, the user selects a specific bastion from a list based on the chosen type
// isPrivateRFC1918 checks if an IP is in RFC1918 ranges.
func isPrivateRFC1918(ip net.IP) bool {
	if ip == nil {
		return false
	}
	// Normalize to 4-byte if IPv4-mapped
	if v4 := ip.To4(); v4 != nil {
		// 10.0.0.0/8
		if v4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if v4[0] == 172 && v4[1] >= 16 && v4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if v4[0] == 192 && v4[1] == 168 {
			return true
		}
	}
	return false
}

// extractHostname attempts to parse a URL or raw host:port and returns hostname.
func extractHostname(endpoint string) string {
	if endpoint == "" {
		return ""
	}
	// If it looks like a URL, parse it
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		if u, err := url.Parse(endpoint); err == nil {
			return u.Hostname()
		}
	}
	// Otherwise, strip optional :port
	host := endpoint
	if i := strings.Index(host, ":"); i > -1 {
		host = host[:i]
	}
	return host
}

// resolveHostToIP resolves a hostname to an IP; prefer RFC1918 private IPv4 if available.
func resolveHostToIP(host string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("empty host")
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.String(), nil
	}
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return "", fmt.Errorf("failed to resolve host %s: %w", host, err)
	}
	// Prefer private IPv4
	for _, ip := range ips {
		if v4 := ip.To4(); v4 != nil && isPrivateRFC1918(v4) {
			return v4.String(), nil
		}
	}
	// Fallback: first IPv4
	for _, ip := range ips {
		if v4 := ip.To4(); v4 != nil {
			return v4.String(), nil
		}
	}
	// Last resort: first IP string
	return ips[0].String(), nil
}

func RunCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Create a new bastion service to interact with bastion resources
	service, err := bastionSvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("error creating Bastion service: %w", err)
	}

	// STEP 1: Type Selection
	typeModel := NewTypeSelectionModel()
	typeProgram := tea.NewProgram(typeModel)

	// Run the type selection TUI and wait for user selection
	typeResult, err := typeProgram.Run()
	if err != nil {
		return fmt.Errorf("error running type selection TUI: %w", err)
	}

	// Process the type selection result
	// If the user didn't make a selection or cancelled, exit early
	typeModel, ok := typeResult.(TypeSelectionModel)
	if !ok || typeModel.Choice == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Check if the user selected the Bastion type
	if typeModel.Choice == TypeBastion {
		util.ShowConstructionAnimation()
		return nil
	}

	// STEP 2: Bastion Selection
	// Gets the appropriate list of bastions based on the selected type
	var bastions []bastionSvc.Bastion

	ctx := context.Background()

	if typeModel.Choice == TypeSession {
		bastions, err = service.List(ctx)
		bastions = slices.DeleteFunc(bastions, func(b bastionSvc.Bastion) bool {
			return b.LifecycleState != bastion.BastionLifecycleStateActive
		})
	} else {
		// For Bastion type, also use dummy bastions for now
		// In a real implementation; this might fetch a different set of bastions
		util.ShowConstructionAnimation()
	}

	// Initialize and show the bastion selection UI with the appropriate list
	bastionModel := NewBastionModel(bastions)
	bastionProgram := tea.NewProgram(bastionModel)

	// Run the bastion selection TUI and wait for user selection
	bastionResult, err := bastionProgram.Run()
	if err != nil {
		return fmt.Errorf("error running bastion selection TUI: %w", err)
	}

	// Process the bastion selection result
	if bastionModel, ok := bastionResult.(BastionModel); ok && bastionModel.Choice != "" {
		// User made a selection - find the selected bastion by ID from our list
		var selected bastionSvc.Bastion
		for _, b := range bastions {
			if b.ID == bastionModel.Choice {
				selected = b
				break
			}
		}

		// STEP 3: Session Type Selection
		// Initialize and show the session type selection UI
		sessionTypeModel := NewSessionTypeModel(selected.ID)
		sessionTypeProgram := tea.NewProgram(sessionTypeModel)

		// Run the session type selection TUI and wait for user selection
		sessionTypeResult, err := sessionTypeProgram.Run()
		if err != nil {
			return fmt.Errorf("error running session type selection TUI: %w", err)
		}

		// Process the session type selection result
		sessionTypeModel, ok := sessionTypeResult.(SessionTypeModel)
		if !ok || sessionTypeModel.Choice == "" {
			fmt.Println("Operation cancelled.")
			return nil
		}

		// Check if the user selected Managed SSH
		if sessionTypeModel.Choice == TypeManagedSSH {
			// For Managed SSH, allow choosing a target type and resource similar to Port-Forwarding
			targetTypeModel := NewTargetTypeModel(selected.ID)
			targetTypeProgram := tea.NewProgram(targetTypeModel)
			// Run target selection
			targetTypeResult, err := targetTypeProgram.Run()
			if err != nil {
				return fmt.Errorf("error running target type selection TUI: %w", err)
			}
			// Process result
			targetTypeModel, ok := targetTypeResult.(TargetTypeModel)
			if !ok || targetTypeModel.Choice == "" {
				fmt.Println("Operation cancelled.")
				return nil
			}
			// Handle targets
			if targetTypeModel.Choice == TargetInstance {
				instService, err := instancessvc.NewService(appCtx)
				if err != nil {
					return fmt.Errorf("error creating Instance service: %w", err)
				}
				instances, _, _, err := instService.List(ctx, 300, 0, true)
				if err != nil {
					return fmt.Errorf("error listing instances: %w", err)
				}

				if len(instances) == 0 {
					fmt.Println("No instances found.")
					return nil
				}
				instModel := NewInstanceListModelFancy(instances)
				instProgram := tea.NewProgram(instModel)
				instResult, err := instProgram.Run()
				if err != nil {
					return fmt.Errorf("error running instance selection TUI: %w", err)
				}
				chosen, ok := instResult.(ResourceListModel)
				if !ok || chosen.Choice() == "" {
					fmt.Println("Operation cancelled.")
					return nil
				}
				var selectedInst instancessvc.Instance
				for _, i := range instances {
					if i.ID == chosen.Choice() {
						selectedInst = i
						break
					}
				}
				reachable, reason := service.CanReach(ctx, selected, selectedInst.VcnID, selectedInst.SubnetID)
				if !reachable {
					fmt.Println("Bastion cannot reach selected instance:", reason)
					return nil
				}

				fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to Instance %s.\n",
					sessionTypeModel.Choice, selected.Name, selected.ID, selectedInst.Name)

				// Connect to the instance via Bastion using Managed SSH (no password prompt expected)
				pubKey, privKey := bastionSvc.DefaultSSHKeyPaths()
				sshUser := "opc"
				sessionID, err := service.EnsureManagedSSHSession(ctx, selected.ID, selectedInst.ID, selectedInst.IP, sshUser, 22, pubKey, 0)
				if err != nil {
					return fmt.Errorf("failed to ensure managed SSH session: %w", err)
				}
				region, _ := appCtx.Provider.Region()
				sshCmd := bastionSvc.BuildManagedSSHCommand(privKey, sessionID, region, selectedInst.IP, sshUser)
				fmt.Printf("\nExecuting: %s\n\n", sshCmd)
				cmd := exec.Command("bash", "-lc", sshCmd)
				cmd.Stdout = appCtx.Stdout
				cmd.Stderr = appCtx.Stderr
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("ssh command failed: %w", err)
				}

				return nil
			}
			if targetTypeModel.Choice == TargetDatabase {
				dbService, err := autonomousdbsvc.NewService(appCtx)
				if err != nil {
					return fmt.Errorf("error creating Database service: %w", err)
				}
				dbs, _, _, err := dbService.List(ctx, 50, 0)
				if err != nil {
					return fmt.Errorf("error listing databases: %w", err)
				}
				if len(dbs) == 0 {
					fmt.Println("No Autonomous Databases found.")
					return nil
				}
				dbModel := NewDBListModelFancy(dbs)
				dbProgram := tea.NewProgram(dbModel)
				dbResult, err := dbProgram.Run()
				if err != nil {
					return fmt.Errorf("error running DB selection TUI: %w", err)
				}
				chosen, ok := dbResult.(ResourceListModel)
				if !ok || chosen.Choice() == "" {
					fmt.Println("Operation cancelled.")
					return nil
				}
				var selectedDB autonomousdbsvc.AutonomousDatabase
				for _, d := range dbs {
					if d.ID == chosen.Choice() {
						selectedDB = d
						break
					}
				}
				// Inform type explicitly as a database instance; reachability cannot be auto-verified
				_, reason := service.CanReach(ctx, selected, "", "")
				fmt.Println("Reachability to DB cannot be automatically verified:", reason)
				fmt.Printf("Selected database instance: %s (ID: %s)\n", selectedDB.Name, selectedDB.ID)
				fmt.Printf("\n---\nPrepared %s session on Bastion %s (ID: %s) to database instance %s.\n",
					sessionTypeModel.Choice, selected.Name, selected.ID, selectedDB.Name)
				return nil
			}

			if targetTypeModel.Choice == TargetOKE {
				okeService, err := okesvc.NewService(appCtx)
				if err != nil {
					return fmt.Errorf("error creating OKE service: %w", err)
				}
				list, _, _, err := okeService.List(ctx, 50, 0)
				if err != nil {
					return fmt.Errorf("error listing OKE clusters: %w", err)
				}
				if len(list) == 0 {
					fmt.Println("No OKE clusters found.")
					return nil
				}
				clusterModel := NewOKEListModelFancy(list)
				clusterProgram := tea.NewProgram(clusterModel)
				clusterResult, err := clusterProgram.Run()
				if err != nil {
					return fmt.Errorf("error running OKE selection TUI: %w", err)
				}
				chosen, ok := clusterResult.(ResourceListModel)
				if !ok || chosen.Choice() == "" {
					fmt.Println("Operation cancelled.")
					return nil
				}
				var selectedCluster okesvc.Cluster
				for _, c := range list {
					if c.ID == chosen.Choice() {
						selectedCluster = c
						break
					}
				}
				reachable, reason := service.CanReach(ctx, selected, selectedCluster.VcnID, "")
				if !reachable {
					fmt.Println("Bastion cannot reach selected OKE node (cluster VCN mismatch):", reason)
					return nil
				}
				// Identify output explicitly as OKE node (node selection not implemented; using cluster scope for reachability)
				fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to OKE node in cluster %s.\n",
					sessionTypeModel.Choice, selected.Name, selected.ID, selectedCluster.Name)
				return nil
			}
			// Fallback
			fmt.Printf("\n---\nPrepared %s Session on Bastion: %s (ID: %s)\nTarget: %s\n",
				sessionTypeModel.Choice, selected.Name, selected.ID, targetTypeModel.Choice)
			return nil
		}
		var clusters []okesvc.Cluster
		// STEP 4: Target Type Selection
		if sessionTypeModel.Choice == TypePortForwarding {
			// Initialize and show the target type selection UI
			targetTypeModel := NewTargetTypeModel(selected.ID)
			targetTypeProgram := tea.NewProgram(targetTypeModel)

			// Run the target type selection TUI and wait for user selection
			targetTypeResult, err := targetTypeProgram.Run()
			if err != nil {
				return fmt.Errorf("error running target type selection TUI: %w", err)
			}

			// Process the target type selection result
			targetTypeModel, ok := targetTypeResult.(TargetTypeModel)
			if !ok || targetTypeModel.Choice == "" {
				fmt.Println("Operation cancelled.")
				return nil
			}

			if targetTypeModel.Choice == TargetOKE {
				okeService, err := okesvc.NewService(appCtx)
				if err != nil {
					return fmt.Errorf("error creating OKE service: %w", err)
				}
				list, _, _, err := okeService.List(ctx, 50, 0)
				if err != nil {
					return fmt.Errorf("error listing OKE clusters: %w", err)
				}
				clusters = list
				if len(clusters) == 0 {
					fmt.Println("No OKE clusters found.")
					return nil
				}
				// Select a specific cluster with a fancy searchable list
				clusterModel := NewOKEListModelFancy(clusters)
				clusterProgram := tea.NewProgram(clusterModel)
				clusterResult, err := clusterProgram.Run()
				if err != nil {
					return fmt.Errorf("error running OKE selection TUI: %w", err)
				}
				chosen, ok := clusterResult.(ResourceListModel)
				if !ok || chosen.Choice() == "" {
					fmt.Println("Operation cancelled.")
					return nil
				}
				// Find the selected cluster
				var selectedCluster okesvc.Cluster
				for _, c := range clusters {
					if c.ID == chosen.Choice() {
						selectedCluster = c
						break
					}
				}
				// Reachability check
				reachable, reason := service.CanReach(ctx, selected, selectedCluster.VcnID, "")
				if !reachable {
					fmt.Println("Bastion cannot reach selected OKE cluster:", reason)
					return nil
				}
				fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to OKE cluster %s.\n",
					sessionTypeModel.Choice, selected.Name, selected.ID, selectedCluster.Name)

				// Prompt for local port to forward to OKE API server (remote 6443)
				localPort, err := util.PromptPort("Enter local port to forward to OKE API (remote 6443)", 6443)
				if err != nil {
					return fmt.Errorf("failed to read port: %w", err)
				}
				pubKey, privKey := bastionSvc.DefaultSSHKeyPaths()
				// Resolve the OKE API endpoint host to a private IP suitable for Bastion.
				// Prefer PrivateEndpoint first; fall back to KubernetesEndpoint only if necessary.
				candidates := []string{}
				if h := extractHostname(selectedCluster.PrivateEndpoint); h != "" {
					candidates = append(candidates, h)
				}
				if h := extractHostname(selectedCluster.KubernetesEndpoint); h != "" {
					candidates = append(candidates, h)
				}
				if len(candidates) == 0 {
					return fmt.Errorf("could not determine OKE API host from endpoint: kube=%q private=%q", selectedCluster.KubernetesEndpoint, selectedCluster.PrivateEndpoint)
				}
				var targetIP string
				var lastErr error
				for _, host := range candidates {
					ip, err := resolveHostToIP(host)
					if err == nil {
						targetIP = ip
						break
					}
					lastErr = err
				}
				if targetIP == "" {
					return fmt.Errorf("failed to resolve OKE API endpoint to a private IP: %v", lastErr)
				}
				// Ensure or create the port forwarding bastion session for the cluster API endpoint IP
				sessionID, err := service.EnsurePortForwardSession(ctx, selected.ID, targetIP, 6443, pubKey, 0)
				if err != nil {
					return fmt.Errorf("failed to ensure port forwarding session: %w", err)
				}
				region, _ := appCtx.Provider.Region()
				logFile := fmt.Sprintf("ssh-tunnel-%d.log", localPort)
				sshCmd := bastionSvc.BuildPortForwardNohupCommand(privKey, sessionID, region, targetIP, localPort, 6443, logFile)
				fmt.Printf("\nStarting background OKE API tunnel: %s\n\n", sshCmd)
				cmd := exec.Command("bash", "-lc", sshCmd)
				cmd.Stdout = appCtx.Stdout
				cmd.Stderr = appCtx.Stderr
				// no stdin required for nohup background
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to start SSH tunnel: %w", err)
				}
				fmt.Printf("SSH tunnel to OKE API started in background. Access: https://127.0.0.1:%d (kube-apiserver)\nLogs: %s\n", localPort, logFile)
				return nil
			}

			// Check if the user selected Database
			if targetTypeModel.Choice == TargetDatabase {
				dbService, err := autonomousdbsvc.NewService(appCtx)
				if err != nil {
					return fmt.Errorf("error creating Database service: %w", err)
				}
				dbs, _, _, err := dbService.List(ctx, 50, 0)
				if err != nil {
					return fmt.Errorf("error listing databases: %w", err)
				}
				if len(dbs) == 0 {
					fmt.Println("No Autonomous Databases found.")
					return nil
				}
				dbModel := NewDBListModelFancy(dbs)
				dbProgram := tea.NewProgram(dbModel)
				dbResult, err := dbProgram.Run()
				if err != nil {
					return fmt.Errorf("error running DB selection TUI: %w", err)
				}
				chosen, ok := dbResult.(ResourceListModel)
				if !ok || chosen.Choice() == "" {
					fmt.Println("Operation cancelled.")
					return nil
				}
				var selectedDB autonomousdbsvc.AutonomousDatabase
				for _, d := range dbs {
					if d.ID == chosen.Choice() {
						selectedDB = d
						break
					}
				}
				// We don't have VCN/Subnet info for DBs here; inform the user
				_, reason := service.CanReach(ctx, selected, "", "")
				fmt.Println("Reachability to DB cannot be automatically verified:", reason)
				fmt.Printf("Selected DB: %s (ID: %s)\n", selectedDB.Name, selectedDB.ID)
				return nil
			}

			if targetTypeModel.Choice == TargetInstance {
				instService, err := instancessvc.NewService(appCtx)
				if err != nil {
					return fmt.Errorf("error creating Instance service: %w", err)
				}
				instances, _, _, err := instService.List(ctx, 300, 0, true)
				if err != nil {
					return fmt.Errorf("error listing instances: %w", err)
				}
				if len(instances) == 0 {
					fmt.Println("No instances found.")
					return nil
				}
				instModel := NewInstanceListModelFancy(instances)
				instProgram := tea.NewProgram(instModel)
				instResult, err := instProgram.Run()
				if err != nil {
					return fmt.Errorf("error running instance selection TUI: %w", err)
				}
				chosen, ok := instResult.(ResourceListModel)
				if !ok || chosen.Choice() == "" {
					fmt.Println("Operation cancelled.")
					return nil
				}
				var selectedInst instancessvc.Instance
				for _, i := range instances {
					if i.ID == chosen.Choice() {
						selectedInst = i
						break
					}
				}
				reachable, reason := service.CanReach(ctx, selected, selectedInst.VcnID, selectedInst.SubnetID)
				if !reachable {
					fmt.Println("Bastion cannot reach selected instance:", reason)
					return nil
				}
				fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to Instance %s.\n",
					sessionTypeModel.Choice, selected.Name, selected.ID, selectedInst.Name)

				// Prompt for local port to forward
				port, err := util.PromptPort("Enter port to forward (local:target)", 6443)
				if err != nil {
					return fmt.Errorf("failed to read port: %w", err)
				}
				pubKey, privKey := bastionSvc.DefaultSSHKeyPaths()
				// Ensure or create the port-forwarding bastion session
				sessionID, err := service.EnsurePortForwardSession(ctx, selected.ID, selectedInst.IP, port, pubKey, 0)
				if err != nil {
					return fmt.Errorf("failed to ensure port forwarding session: %w", err)
				}
				region, _ := appCtx.Provider.Region()
				logFile := fmt.Sprintf("ssh-tunnel-%d.log", port)
				sshCmd := bastionSvc.BuildPortForwardNohupCommand(privKey, sessionID, region, selectedInst.IP, port, port, logFile)
				fmt.Printf("\nStarting background tunnel: %s\n\n", sshCmd)
				cmd := exec.Command("bash", "-lc", sshCmd)
				cmd.Stdout = appCtx.Stdout
				cmd.Stderr = appCtx.Stderr
				// no stdin required for nohup background
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to start SSH tunnel: %w", err)
				}
				fmt.Printf("SSH tunnel started in background. Logs: %s\n", logFile)
				return nil
			}

			// Fallback
			fmt.Printf("\n---\nCreated %s Session on Bastion: %s\nID: %s\nTarget: %s\n",
				sessionTypeModel.Choice, selected.Name, selected.ID, targetTypeModel.Choice)
		}
	} else {
		fmt.Println("Operation cancelled.")
	}

	return nil
}
