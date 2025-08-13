package bastion

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"

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

// handleTypeSelection handles the selection of the bastion type
func handleTypeSelection() (BastionType, error) {
	// Initialize and show the type selection UI
	typeModel := NewTypeSelectionModel()
	typeProgram := tea.NewProgram(typeModel)

	// Run the type selection TUI and wait for user selection
	typeResult, err := typeProgram.Run()
	if err != nil {
		return "", fmt.Errorf("error running type selection TUI: %w", err)
	}

	// Process the type selection result
	typeModel, ok := typeResult.(TypeSelectionModel)
	if !ok || typeModel.Choice == "" {
		return "", nil // User cancelled
	}

	return typeModel.Choice, nil
}

// handleBastionSelection handles the selection of a bastion from a list
func handleBastionSelection(ctx context.Context, service *bastionSvc.Service, typeChoice BastionType) (bastionSvc.Bastion, error) {
	var bastions []bastionSvc.Bastion
	var err error

	// Get the appropriate list of bastions based on the selected type
	if typeChoice == TypeSession {
		bastions, err = service.List(ctx)
		if err != nil {
			return bastionSvc.Bastion{}, fmt.Errorf("error listing bastions: %w", err)
		}

		// Filter out non-active bastions
		bastions = slices.DeleteFunc(bastions, func(b bastionSvc.Bastion) bool {
			return b.LifecycleState != bastion.BastionLifecycleStateActive
		})
	} else {
		// For other types, show construction animation
		util.ShowConstructionAnimation()
		return bastionSvc.Bastion{}, nil
	}

	// Initialize and show the bastion selection UI
	bastionModel := NewBastionModel(bastions)
	bastionProgram := tea.NewProgram(bastionModel)

	// Run the bastion selection TUI and wait for user selection
	bastionResult, err := bastionProgram.Run()
	if err != nil {
		return bastionSvc.Bastion{}, fmt.Errorf("error running bastion selection TUI: %w", err)
	}

	// Process the bastion selection result
	bastionModel, ok := bastionResult.(BastionModel)
	if !ok || bastionModel.Choice == "" {
		return bastionSvc.Bastion{}, nil // User cancelled
	}

	// Find the selected bastion by ID
	var selected bastionSvc.Bastion
	for _, b := range bastions {
		if b.ID == bastionModel.Choice {
			selected = b
			break
		}
	}

	return selected, nil
}

// handleSessionTypeSelection handles the selection of the session type
func handleSessionTypeSelection(bastionID string) (SessionType, error) {
	// Initialize and show the session type selection UI
	sessionTypeModel := NewSessionTypeModel(bastionID)
	sessionTypeProgram := tea.NewProgram(sessionTypeModel)

	// Run the session type selection TUI and wait for user selection
	sessionTypeResult, err := sessionTypeProgram.Run()
	if err != nil {
		return "", fmt.Errorf("error running session type selection TUI: %w", err)
	}

	// Process the session type selection result
	sessionTypeModel, ok := sessionTypeResult.(SessionTypeModel)
	if !ok || sessionTypeModel.Choice == "" {
		return "", nil // User cancelled
	}

	return sessionTypeModel.Choice, nil
}

// handleTargetTypeSelection handles the selection of the target type
func handleTargetTypeSelection(bastionID string) (TargetType, error) {
	// Initialize and show the target type selection UI
	targetTypeModel := NewTargetTypeModel(bastionID)
	targetTypeProgram := tea.NewProgram(targetTypeModel)

	// Run the target type selection TUI and wait for user selection
	targetTypeResult, err := targetTypeProgram.Run()
	if err != nil {
		return "", fmt.Errorf("error running target type selection TUI: %w", err)
	}

	// Process the target type selection result
	targetTypeModel, ok := targetTypeResult.(TargetTypeModel)
	if !ok || targetTypeModel.Choice == "" {
		return "", nil // User cancelled
	}

	return targetTypeModel.Choice, nil
}

// handleInstanceTarget handles the selection and connection to an instance target
func handleInstanceTarget(ctx context.Context, appCtx *app.ApplicationContext, service *bastionSvc.Service,
	selected bastionSvc.Bastion, sessionType SessionType) error {

	// Create an instance service
	instService, err := instancessvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("error creating Instance service: %w", err)
	}

	// List instances
	instances, _, _, err := instService.List(ctx, 300, 0, true)
	if err != nil {
		return fmt.Errorf("error listing instances: %w", err)
	}

	if len(instances) == 0 {
		fmt.Println("No instances found.")
		return nil
	}

	// Show instance selection UI
	instModel := NewInstanceListModelFancy(instances)
	instProgram := tea.NewProgram(instModel)
	instResult, err := instProgram.Run()
	if err != nil {
		return fmt.Errorf("error running instance selection TUI: %w", err)
	}

	// Process selection result
	chosen, ok := instResult.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Find a selected instance
	var selectedInst instancessvc.Instance
	for _, i := range instances {
		if i.ID == chosen.Choice() {
			selectedInst = i
			break
		}
	}

	// Check if bastion can reach the instance
	reachable, reason := service.CanReach(ctx, selected, selectedInst.VcnID, selectedInst.SubnetID)
	if !reachable {
		fmt.Println("Bastion cannot reach selected instance:", reason)
		return nil
	}

	fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to Instance %s.\n",
		sessionType, selected.Name, selected.ID, selectedInst.Name)

	// Handle different session types
	if sessionType == TypeManagedSSH {
		// Connect to the instance via Bastion using Managed SSH
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
	} else if sessionType == TypePortForwarding {
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
	}

	return nil
}

// handleDatabaseTarget handles the selection and connection to a database target
func handleDatabaseTarget(ctx context.Context, appCtx *app.ApplicationContext, service *bastionSvc.Service,
	selected bastionSvc.Bastion, sessionType SessionType) error {

	// Create a database service
	dbService, err := autonomousdbsvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("error creating Database service: %w", err)
	}

	// List databases
	dbs, _, _, err := dbService.List(ctx, 50, 0)
	if err != nil {
		return fmt.Errorf("error listing databases: %w", err)
	}

	if len(dbs) == 0 {
		fmt.Println("No Autonomous Databases found.")
		return nil
	}

	// Show database selection UI
	dbModel := NewDBListModelFancy(dbs)
	dbProgram := tea.NewProgram(dbModel)
	dbResult, err := dbProgram.Run()
	if err != nil {
		return fmt.Errorf("error running DB selection TUI: %w", err)
	}

	// Process selection result
	chosen, ok := dbResult.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Find the selected database
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
		sessionType, selected.Name, selected.ID, selectedDB.Name)

	return nil
}

// handleOKETarget handles the selection and connection to an OKE cluster target
func handleOKETarget(ctx context.Context, appCtx *app.ApplicationContext, service *bastionSvc.Service,
	selected bastionSvc.Bastion, sessionType SessionType) error {

	// Create an OKE service
	okeService, err := okesvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("error creating OKE service: %w", err)
	}

	// List OKE clusters
	clusters, _, _, err := okeService.List(ctx, 50, 0)
	if err != nil {
		return fmt.Errorf("error listing OKE clusters: %w", err)
	}

	if len(clusters) == 0 {
		fmt.Println("No OKE clusters found.")
		return nil
	}

	// Show cluster selection UI
	clusterModel := NewOKEListModelFancy(clusters)
	clusterProgram := tea.NewProgram(clusterModel)
	clusterResult, err := clusterProgram.Run()
	if err != nil {
		return fmt.Errorf("error running OKE selection TUI: %w", err)
	}

	// Process selection result
	chosen, ok := clusterResult.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Find selected cluster
	var selectedCluster okesvc.Cluster
	for _, c := range clusters {
		if c.ID == chosen.Choice() {
			selectedCluster = c
			break
		}
	}

	// Check if bastion can reach the cluster
	reachable, reason := service.CanReach(ctx, selected, selectedCluster.VcnID, "")
	if !reachable {
		fmt.Println("Bastion cannot reach selected OKE cluster:", reason)
		return nil
	}

	fmt.Printf("\n---\nValidated %s session on Bastion %s (ID: %s) to OKE cluster %s.\n",
		sessionType, selected.Name, selected.ID, selectedCluster.Name)

	// Handle different session types
	if sessionType == TypePortForwarding {
		// Prompt for local port to forward to OKE API server (remote 6443)
		localPort, err := util.PromptPort("Enter local port to forward to OKE API (remote 6443)", 6443)
		if err != nil {
			return fmt.Errorf("failed to read port: %w", err)
		}

		pubKey, privKey := bastionSvc.DefaultSSHKeyPaths()

		// Resolve the OKE API endpoint host to a private IP suitable for Bastion
		candidates := []string{}
		if h := extractHostname(selectedCluster.PrivateEndpoint); h != "" {
			candidates = append(candidates, h)
		}
		if h := extractHostname(selectedCluster.KubernetesEndpoint); h != "" {
			candidates = append(candidates, h)
		}

		if len(candidates) == 0 {
			return fmt.Errorf("could not determine OKE API host from endpoint: kube=%q private=%q",
				selectedCluster.KubernetesEndpoint, selectedCluster.PrivateEndpoint)
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

		fmt.Printf("SSH tunnel to OKE API started in background. Access: https://127.0.0.1:%d (kube-apiserver)\nLogs: %s\n",
			localPort, logFile)
	}

	return nil
}

// RunCreateCommand orchestrates the creation of a bastion session
func RunCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Create a new bastion service to interact with bastion resources
	service, err := bastionSvc.NewService(appCtx)
	if err != nil {
		return fmt.Errorf("error creating Bastion service: %w", err)
	}

	ctx := context.Background()

	// STEP 1: Type Selection
	typeChoice, err := handleTypeSelection()
	if err != nil {
		return err
	}
	if typeChoice == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Check if the user selected the Bastion type
	if typeChoice == TypeBastion {
		util.ShowConstructionAnimation()
		return nil
	}

	// STEP 2: Bastion Selection
	selected, err := handleBastionSelection(ctx, service, typeChoice)
	if err != nil {
		return err
	}
	if selected.ID == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// STEP 3: Session Type Selection
	sessionType, err := handleSessionTypeSelection(selected.ID)
	if err != nil {
		return err
	}
	if sessionType == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// STEP 4: Target Type Selection
	targetType, err := handleTargetTypeSelection(selected.ID)
	if err != nil {
		return err
	}
	if targetType == "" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// STEP 5: Handle a specific target type based on a session type
	switch targetType {
	case TargetInstance:
		return handleInstanceTarget(ctx, appCtx, service, selected, sessionType)
	case TargetDatabase:
		return handleDatabaseTarget(ctx, appCtx, service, selected, sessionType)
	case TargetOKE:
		return handleOKETarget(ctx, appCtx, service, selected, sessionType)
	default:
		// Fallback for unknown target types
		fmt.Printf("\n---\nPrepared %s Session on Bastion: %s (ID: %s)\nTarget: %s\n",
			sessionType, selected.Name, selected.ID, targetType)
		return nil
	}
}
