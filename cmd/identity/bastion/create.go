package bastion

import (
	"context"
	"fmt"
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

// RunCreateCommand executes the create command with an interactive TUI
// It implements a two-step selection process:
// 1. First, the user selects a type (Bastion or Session)
// 2. Then, the user selects a specific bastion from a list based on the chosen type
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
				instances, _, _, err := instService.List(ctx, 50, 0, true)
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
				instances, _, _, err := instService.List(ctx, 50, 0, true)
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
