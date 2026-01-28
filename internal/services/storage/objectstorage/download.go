package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	osadapter "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// DownloadFile handles the interactive download flow:
// 1. Show bucket list TUI to select the source bucket
// 2. Show object list TUI to select object to download
// 3. Download the file with progress TUI
func DownloadFile(appCtx *app.ApplicationContext) error {
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

	// Step 1: Show bucket list TUI to select the source bucket
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

	// Step 2: Show object list TUI to select an object
	objects, err := service.ListObjects(ctx, namespace, bucketName)
	if err != nil {
		return fmt.Errorf("listing objects: %w", err)
	}

	if len(objects) == 0 {
		fmt.Fprintf(appCtx.Stdout, "Bucket %q is empty.\n", bucketName)
		return nil
	}

	objectModel := osadapter.NewObjectListModel(objects, bucketName)
	objectName, err := tui.Run(objectModel)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting object: %w", err)
	}

	// Step 3: Download the file with progress TUI
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	title := fmt.Sprintf("Downloading %s from %s", objectName, bucketName)
	progressRunner := tui.NewProgressRunner(title)
	progressRunner.Start()

	// Channel to signal completion
	done := make(chan error, 1)

	// Start download in goroutine
	go func() {
		progressFn := func(p TransferProgress) {
			percent := float64(p.BytesTransferred) / float64(p.TotalBytes)
			bytesInfo := fmt.Sprintf("%s / %s",
				util.HumanizeBytesIEC(p.BytesTransferred),
				util.HumanizeBytesIEC(p.TotalBytes))

			status := "Downloading..."
			if percent >= 1.0 {
				status = "Complete!"
			}

			progressRunner.UpdateProgress(percent, bytesInfo, "", status)
		}

		err := service.DownloadObject(ctx, namespace, bucketName, objectName, cwd, progressFn)
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

	// Wait for download to finish and check a result
	downloadErr := <-done
	if downloadErr != nil {
		return fmt.Errorf("downloading object: %w", downloadErr)
	}

	fmt.Fprintf(appCtx.Stdout, "\nSuccessfully downloaded %s to %s\n", objectName, cwd)
	return nil
}
