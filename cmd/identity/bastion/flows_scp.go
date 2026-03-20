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

// connectSCP runs the SCP flow for an Instance target:
// 1. Select instance
// 2. Select SSH key pair
// 3. Select local file to upload (file picker TUI)
// 4. Prompt for remote destination path
// 5. Create managed SSH session
// 6. Execute SCP with progress
func connectSCP(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
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

	// Step 3: Select local file to upload using the file picker TUI
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	fileModel, err := tui.NewFilePickerModel(cwd)
	if err != nil {
		return fmt.Errorf("creating file picker: %w", err)
	}

	localFile, err := tui.RunFilePicker(fileModel)
	if err != nil {
		return ErrAborted
	}

	// Step 4: Prompt for remote destination path
	remotePath, err := util.PromptString("Enter remote destination path", "/tmp/"+filepath.Base(localFile))
	if err != nil {
		return fmt.Errorf("read remote path: %w", err)
	}

	// Step 5: Prompt for SSH username
	sshUser, err := util.PromptString("Enter SSH username", "opc")
	if err != nil {
		return fmt.Errorf("read ssh username: %w", err)
	}

	logger.Logger.Info("Creating SCP session",
		"bastion", b.DisplayName,
		"instance", inst.DisplayName,
		"file", filepath.Base(localFile),
		"remote_path", remotePath,
	)

	// Step 6: Create managed SSH session (SCP uses managed SSH under the hood)
	sessID, err := svc.EnsureManagedSSHSession(ctx, b.OCID, inst.OCID, inst.PrimaryIP, sshUser, 22, pubKey, 0)
	if err != nil {
		return fmt.Errorf("ensure managed SSH session for SCP: %w", err)
	}

	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	// Step 7: Build and execute SCP command with progress TUI
	scpCmd := bastionSvc.BuildSCPCommand(privKey, sessID, region, inst.PrimaryIP, sshUser, localFile, remotePath)

	title := fmt.Sprintf("Uploading %s to %s:%s", filepath.Base(localFile), inst.DisplayName, remotePath)
	progressRunner := tui.NewProgressRunner(title)
	progressRunner.Start()

	done := make(chan error, 1)

	go func() {
		fileInfo, statErr := os.Stat(localFile)
		if statErr != nil {
			progressRunner.SendError(statErr)
			done <- statErr
			return
		}

		totalBytes := fileInfo.Size()
		bytesInfo := fmt.Sprintf("0 B / %s", util.HumanizeBytesIEC(totalBytes))
		progressRunner.UpdateProgress(0, bytesInfo, "", "Connecting via bastion...")

		progressRunner.UpdateProgress(0.05, bytesInfo, "", "SCP transfer in progress...")

		err := bastionSvc.RunShell(ctx, appCtx.Stdout, appCtx.Stderr, scpCmd)
		if err != nil {
			progressRunner.SendError(err)
			done <- err
			return
		}

		bytesInfo = fmt.Sprintf("%s / %s", util.HumanizeBytesIEC(totalBytes), util.HumanizeBytesIEC(totalBytes))
		progressRunner.UpdateProgress(1.0, bytesInfo, "", "Complete!")
		done <- nil
	}()

	if err := progressRunner.Run(); err != nil {
		return err
	}

	uploadErr := <-done
	if uploadErr != nil {
		return fmt.Errorf("SCP transfer: %w", uploadErr)
	}

	logger.Logger.Info("Successfully copied file to remote instance",
		"file", filepath.Base(localFile),
		"instance", inst.DisplayName,
		"remote_path", remotePath,
	)
	return nil
}
