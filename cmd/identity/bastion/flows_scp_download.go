package bastion

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociInst "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
	instSvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// connectSCPDownload runs the SCP download flow for an Instance target:
// 1. Select instance
// 2. Select SSH key pair
// 3. Enter remote file/directory path (TUI text input)
// 4. Enter local destination path (TUI text input)
// 5. Create managed SSH session
// 6. Execute SCP download with progress
func connectSCPDownload(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion) error {

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

	if len(instances) == 0 {
		logger.Logger.Info("No instances found.")
		return nil
	}

	// Step 1: Select instance
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

	var inst instSvc.Instance
	for _, it := range instances {
		if it.OCID == chosen.Choice() {
			inst = it
			break
		}
	}

	// Step 2: Select SSH key pair
	pubKey, privKey, err := SelectSSHKeyPair(ctx)
	if err != nil {
		return err
	}

	// Validate bastion can reach instance
	if ok, reason := svc.CanReach(ctx, b, inst.VcnID, inst.SubnetID); !ok {
		logger.Logger.Info("Bastion cannot reach selected instance", "reason", reason)
		return nil
	}

	// Step 3: Enter remote file/directory path to download (TUI text input)
	remoteInput := tui.NewTextInputModel(
		"Enter remote file or directory path to download",
		"/home/opc/example.log",
		"",
	)
	remotePath, err := tui.RunTextInput(remoteInput)
	if err != nil {
		return ErrAborted
	}
	if remotePath == "" {
		logger.Logger.Info("No remote path specified, aborting.")
		return nil
	}

	// Step 4: Enter local destination path (TUI text input)
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	defaultLocal := filepath.Join(cwd, filepath.Base(remotePath))
	localInput := tui.NewTextInputModel(
		"Enter local destination path",
		defaultLocal,
		defaultLocal,
	)
	localPath, err := tui.RunTextInput(localInput)
	if err != nil {
		return ErrAborted
	}
	if localPath == "" {
		localPath = defaultLocal
	}

	// Step 5: Prompt for SSH username
	sshUser, err := util.PromptString("Enter SSH username", "opc")
	if err != nil {
		return fmt.Errorf("read ssh username: %w", err)
	}

	logger.Logger.Info("Creating SCP download session",
		"bastion", b.DisplayName,
		"instance", inst.DisplayName,
		"remote_path", remotePath,
		"local_path", localPath,
	)

	// Step 6: Create managed SSH session (SCP uses managed SSH under the hood)
	sessID, err := svc.EnsureManagedSSHSession(ctx, b.OCID, inst.OCID, inst.PrimaryIP, sshUser, 22, pubKey, 0)
	if err != nil {
		return fmt.Errorf("ensure managed SSH session for SCP download: %w", err)
	}

	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	// Step 7: Build and execute SCP download command with progress TUI
	scpCmd := bastionSvc.BuildSCPDownloadCommand(privKey, sessID, region, inst.PrimaryIP, sshUser, remotePath, localPath)

	title := fmt.Sprintf("Downloading %s:%s to %s", inst.DisplayName, remotePath, localPath)
	progressRunner := tui.NewProgressRunner(title)
	progressRunner.Start()

	done := make(chan error, 1)

	go func() {
		progressRunner.UpdateProgress(0.05, "", "", "Connecting via bastion...")

		err := bastionSvc.RunShell(ctx, appCtx.Stdout, appCtx.Stderr, scpCmd)
		if err != nil {
			progressRunner.SendError(err)
			done <- err
			return
		}

		progressRunner.UpdateProgress(1.0, "", "", "Complete!")
		done <- nil
	}()

	if err := progressRunner.Run(); err != nil {
		return err
	}

	downloadErr := <-done
	if downloadErr != nil {
		return fmt.Errorf("SCP download: %w", downloadErr)
	}

	logger.Logger.Info("Successfully downloaded from remote instance",
		"remote_path", remotePath,
		"instance", inst.DisplayName,
		"local_path", localPath,
	)
	return nil
}
