package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	osadapter "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// UploadFile handles the interactive upload flow:
// 1. Show bucket list TUI to select destination bucket
// 2. Show file picker TUI to select file to upload
// 3. Upload the file with progress TUI
func UploadFile(appCtx *app.ApplicationContext) error {
	ctx := context.Background()
	client, err := oci.NewObjectStorageClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating object storage client: %w", err)
	}

	bucketAdapter := osadapter.NewAdapter(client)
	service := NewService(bucketAdapter, appCtx.Logger, appCtx.CompartmentID)

	// Get namespace
	namespace, err := service.GetNamespace(ctx)
	if err != nil {
		return fmt.Errorf("getting namespace: %w", err)
	}

	// Step 1: Show bucket list TUI to select destination bucket
	buckets, err := service.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("listing buckets: %w", err)
	}

	if len(buckets) == 0 {
		fmt.Fprintln(appCtx.Stdout, "No buckets found in the compartment.")
		return nil
	}

	bucketModel := osadapter.NewBucketListModel(buckets)
	bucketID, err := tui.Run(bucketModel)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting bucket: %w", err)
	}

	// Get bucket name from OCID
	bucketName, err := service.osRepo.GetBucketNameByOCID(ctx, appCtx.CompartmentID, bucketID)
	if err != nil {
		return fmt.Errorf("getting bucket name: %w", err)
	}

	// Step 2: Show file picker TUI to select file
	fileModel, err := tui.NewFilePickerModel(".")
	if err != nil {
		return fmt.Errorf("creating file picker: %w", err)
	}

	filePath, err := tui.RunFilePicker(fileModel)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting file: %w", err)
	}

	// Get the object name (just the filename, not the full path)
	objectName := filepath.Base(filePath)

	// Step 3: Upload the file with progress TUI
	title := fmt.Sprintf("Uploading %s to %s", objectName, bucketName)
	progressRunner := tui.NewProgressRunner(title)
	progressRunner.Start()

	// Channel to signal completion
	done := make(chan error, 1)

	// Start upload in goroutine
	go func() {
		progressFn := func(p TransferProgress) {
			percent := float64(p.BytesTransferred) / float64(p.TotalBytes)
			bytesInfo := fmt.Sprintf("%s / %s",
				util.HumanizeBytesIEC(p.BytesTransferred),
				util.HumanizeBytesIEC(p.TotalBytes))

			extraInfo := ""
			if p.TotalParts > 1 {
				extraInfo = fmt.Sprintf("Part %d/%d", p.PartNumber, p.TotalParts)
			}

			status := "Uploading..."
			if percent >= 1.0 {
				status = "Complete!"
			}

			progressRunner.UpdateProgress(percent, bytesInfo, extraInfo, status)
		}

		err := service.UploadObject(ctx, namespace, bucketName, objectName, filePath, progressFn)
		if err != nil {
			progressRunner.SendError(err)
			done <- err
			return
		}
		progressRunner.SendDone()
		done <- nil
	}()

	// Run the progress TUI (blocks until complete)
	if err := progressRunner.Run(); err != nil {
		return err
	}

	// Wait for upload to finish and check result
	uploadErr := <-done
	if uploadErr != nil {
		return fmt.Errorf("uploading object: %w", uploadErr)
	}

	fmt.Fprintf(appCtx.Stdout, "\nSuccessfully uploaded %s to bucket %s\n", objectName, bucketName)
	return nil
}
